package article

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/jarmo/backpocket/template"

	readability "github.com/go-shiori/go-readability"
)

func Create(params ArticleParams) string {
	resp, err := http.Get(params.Url.String())
	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			if content, err := ioutil.ReadAll(resp.Body); err == nil {
				contentType := resp.Header.Get("Content-Type")
				if strings.Contains(contentType, "text/html") {
					article, err := readability.FromReader(bytes.NewReader(content), params.Url.String())
					if err == nil {
						return createReadableArticle(params, article)
					} else {
						return createNonReadableArticle(params, err)
					}
				} else {
					return createNonHTMLContent(params, contentType, content)
				}
			} else {
				panic(err)
			}
		} else {
			return createNonReadableArticle(params, err)
		}
	} else {
		return createNonReadableArticle(params, err)
	}
}

func createReadableArticle(params ArticleParams, article readability.Article) string {
	articleFilePath := ReadableArticleFilePath(params, article)
	articleFile, _ := os.Create(articleFilePath)
	defer articleFile.Close()
	articleFile.WriteString(template.Render(template.ReadableArticleHTML(), template.CreateReadableArticleRenderArgs(params.Url, params.ArchivedAt, article)))
	return articleFilePath
}

func createNonReadableArticle(params ArticleParams, err error) string {
	articleFilePath := NonReadableArticleFilePath(params)
	articleFile, _ := os.Create(articleFilePath)
	defer articleFile.Close()
	articleFile.WriteString(template.Render(template.NonReadableArticleHTML(), template.CreateNonReadableArticleRenderArgs(params.Url, params.ArchivedAt, err)))
	return articleFilePath
}

func createNonHTMLContent(params ArticleParams, contentType string, content []byte) string {
	contentFilePath := NonHTMLContentFilePath(params, contentType)
	contentFile, _ := os.Create(contentFilePath)
	defer contentFile.Close()
	contentFile.Write(content)
	return contentFilePath
}
