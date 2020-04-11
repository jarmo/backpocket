package template

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/jarmo/backpocket/htmlparser"
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

	doc, _ := html.Parse(strings.NewReader(bufferString.String()))

	return htmlparser.Render(
		contentWithAbsoluteIframeUrls(
			args.Address.Scheme,
			contentWithBase64DataSourceImages(doc)))
}

func contentWithAbsoluteIframeUrls(articleScheme string, rootNode *html.Node) *html.Node {
	htmlparser.ForEachNode(rootNode, func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "iframe" {
			srcValue := htmlparser.AttrByName(node, "src")
			if strings.HasPrefix(srcValue, "//") {
				nodeParent := node.Parent
				var attributes []html.Attribute
				for _, attr := range node.Attr {
					if attr.Key == "src" {
						attributes = append(attributes, html.Attribute{Key: "src", Val: articleScheme + ":" + srcValue})
					} else {
						attributes = append(attributes, attr)
					}
				}
				newNode := &html.Node{
					Type: html.ElementNode,
					Data: "iframe",
					Attr: attributes}
				nodeParent.InsertBefore(newNode, node)
				nodeParent.RemoveChild(node)
			}
		}
	})

	return rootNode
}

func contentWithBase64DataSourceImages(rootNode *html.Node) *html.Node {
	htmlparser.ForEachNode(rootNode, func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "img" {
			if srcSetValue := htmlparser.AttrByName(node, "srcset"); srcSetValue != "" {
				replaceImageWithBase64DataSource(node, bestImageSrcSetValue(srcSetValue))
			} else {
				imageSource, _ := url.Parse(htmlparser.AttrByName(node, "src"))
				replaceImageWithBase64DataSource(node, imageSource)
			}
		}
	})

	return rootNode
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
