package template

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

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
							resp, err := http.Get(imageSource.String())
							defer resp.Body.Close()

							if err == nil {
								if imageBytes, err := ioutil.ReadAll(resp.Body); err == nil {
									contentType := http.DetectContentType(imageBytes)
									base64Image := base64.StdEncoding.EncodeToString(imageBytes)
									doc = strings.ReplaceAll(doc, imageSource.String(), fmt.Sprintf("data:%s;base64,%s", contentType, base64Image))
								}
							}
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
