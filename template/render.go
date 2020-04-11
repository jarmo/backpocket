package template

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/url"
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

	return contentWithBase64DataSourceImages(bufferString.String())
}

func contentWithBase64DataSourceImages(doc string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(doc))
	for {
		if tokenizer.Next() == html.ErrorToken {
			break
		}

		if tagName, _ := tokenizer.TagName(); string(tagName) == "img" {
			for {
				attrName, attrValue, hasMoreAttrs := tokenizer.TagAttr()
				if string(attrName) == "src" {
					if imageSource, err := url.Parse(string(attrValue)); err == nil {
						if imageSource.Scheme == "https" || imageSource.Scheme == "http" {
							doc = replaceImageWithBase64DataSource(doc, imageSource)
						}
					}
				}
				if !hasMoreAttrs {
					break
				}
			}
		}
	}
	return doc
}

func replaceImageWithBase64DataSource(doc string, imageSource *url.URL) string {
	resp, err := http.Get(imageSource)

	if err == nil {
		defer resp.Body.Close()
		if imageBytes, err := ioutil.ReadAll(resp.Body); err == nil {
			base64Image := base64.StdEncoding.EncodeToString(imageBytes)
			return strings.ReplaceAll(doc, imageSource.String(), fmt.Sprintf("data:%s;base64,%s", resp.Header.Get("Content-Type"), base64Image))
		}
	}

	return doc
}

