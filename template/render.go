package template

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/jarmo/backpocket/http"

	"golang.org/x/net/html"
)

func Render(content string, args RenderArgs) string {
	tmpl, err := template.New("html").Parse(content)
	if err != nil {
		panic(err)
	}

	bufferString := bytes.NewBufferString("")
	err = tmpl.Execute(bufferString, args)
	if err != nil {
		panic(err)
	}

	return renderNode(contentWithBase64DataSourceImages(bufferString.String()))
}

func renderNode(node *html.Node) string {
	var buf bytes.Buffer
	html.Render(io.Writer(&buf), node)
	return buf.String()
}

func contentWithBase64DataSourceImages(content string) *html.Node {
	doc, _ := html.Parse(strings.NewReader(content))

	forEachNode(doc, func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "img" {
			if srcSetValue := attrByName(node, "srcset"); srcSetValue != "" {
				replaceImageWithBase64DataSource(node, bestImageSrcSetValue(srcSetValue))
			} else {
				imageSource, _ := url.Parse(attrByName(node, "src"))
				replaceImageWithBase64DataSource(node, imageSource)
			}
		}
	})

	return doc
}

func bestImageSrcSetValue(srcSetValue string) *url.URL {
	imageSources := strings.Split(srcSetValue, ",")
	var bestImageSource string
	var bestImageSourceSize = 0
	replaceNonNumericCharacters := regexp.MustCompile("[^0-9]")
	for _, imageSource := range imageSources {
		imageSourceParts := strings.Split(strings.TrimSpace(imageSource), " ")
		imageSizeAsString := replaceNonNumericCharacters.ReplaceAllString(string(imageSourceParts[1]), "")
		if imageSize, err := strconv.Atoi(imageSizeAsString); err == nil {
			if imageSize > bestImageSourceSize {
				bestImageSourceSize = imageSize
				bestImageSource = string(imageSourceParts[0])
			}
		}
	}

	bestImageSourceUrl, _ := url.Parse(bestImageSource)
	return bestImageSourceUrl
}

type nodeFn func(*html.Node)

func forEachNode(rootNode *html.Node, fn nodeFn) {
	var walker func(*html.Node)
	walker = func(node *html.Node) {
		fn(node)
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walker(child)
		}
	}
	walker(rootNode)
}

func replaceImageWithBase64DataSource(node *html.Node, imageSource *url.URL) {
	nodeParent := node.Parent

	if base64DataSource, err := imageAsBase64DataSource(imageSource); err == nil {
		attributes := []html.Attribute{
			html.Attribute{Key: "src", Val: base64DataSource}}
		newNode := &html.Node{
			Type: html.ElementNode,
			Data: "img",
			Attr: attributes}
		nodeParent.InsertBefore(newNode, node)
	}

	nodeParent.RemoveChild(node)
}

func attrByName(node *html.Node, name string) string {
	for _, attr := range node.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}

	return ""
}

func imageAsBase64DataSource(imageSource *url.URL) (string, error) {
	if imageSource.Scheme != "https" && imageSource.Scheme != "http" {
		return "", errors.New("Not supported scheme!")
	}

	if resp, err := http.Get(imageSource); err == nil {
		defer resp.Body.Close()
		if imageBytes, err := ioutil.ReadAll(resp.Body); err == nil {
			base64Image := base64.StdEncoding.EncodeToString(imageBytes)
			return fmt.Sprintf("data:%s;base64,%s", resp.Header.Get("Content-Type"), base64Image), nil
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}
