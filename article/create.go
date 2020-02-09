package article

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/jarmo/backpocket/template"

	readability "github.com/go-shiori/go-readability"
)

func Create(url *url.URL) string {
	resp, err := http.Get(url.String())
	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			if content, err := ioutil.ReadAll(resp.Body); err == nil {
				contentType := http.DetectContentType(content)
				if strings.Contains(contentType, "text/html") {
					article, err := readability.FromReader(bytes.NewReader(content), url.String())
					if err == nil {
						return createReadableArticle(url, article)
					} else {
						return createNonReadableArticle(url, err)
					}
				} else {
					return createNonHTMLContent(url, contentType, content)
				}
			} else {
				panic(err)
			}
		} else {
			return createNonReadableArticle(url, err)
		}
	} else {
		return createNonReadableArticle(url, err)
	}
}

func createReadableArticle(url *url.URL, article readability.Article) string {
	articleFilePath := ReadableArticleFilePath(url, article)
	articleFile, _ := os.Create(articleFilePath)
	defer articleFile.Close()
	articleFile.WriteString(template.Render(template.ReadableArticleHTML(), template.CreateReadableArticleRenderArgs(url, article)))
	return articleFilePath
}

func createNonReadableArticle(url *url.URL, err error) string {
	articleFilePath := NonReadableArticleFilePath(url)
	articleFile, _ := os.Create(articleFilePath)
	defer articleFile.Close()
	articleFile.WriteString(template.Render(template.NonReadableArticleHTML(), template.CreateNonReadableArticleRenderArgs(url, err)))
	return articleFilePath
}

func createNonHTMLContent(url *url.URL, contentType string, content []byte) string {
	contentFilePath := NonHTMLContentFilePath(url, contentType)
	contentFile, _ := os.Create(contentFilePath)
	defer contentFile.Close()
	contentFile.Write(content)
	return contentFilePath
}
