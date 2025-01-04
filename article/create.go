package article

import (
	"bytes"
	"io/ioutil"
	goHttp "net/http"
	"os"
	"strings"

	"github.com/jarmo/backpocket/template"
	"github.com/jarmo/backpocket/http"

	readability "github.com/go-shiori/go-readability"
)

func Create(storageDir string, params ArticleParams) string {
	resp, err := http.Get(params.Url)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == goHttp.StatusOK {
			if content, err := ioutil.ReadAll(resp.Body); err == nil {
				contentType := resp.Header.Get("Content-Type")
				if strings.Contains(contentType, "text/html") {
					metaRefreshUrl := HttpEquivRefreshUrl(params.Url, string(content[:]))
					if metaRefreshUrl != nil {
						return Create(storageDir, ArticleParams{Url: metaRefreshUrl, ArchivedAt: params.ArchivedAt})
					} else {
						article, err := readability.FromReader(bytes.NewReader(content), params.Url)
						if err == nil {
							return createReadableArticle(storageDir, params, article)
						} else {
							return createNonReadableArticle(storageDir, params, err)
						}
					}
				} else {
					return createNonHTMLContent(storageDir, params, contentType, content)
				}
			} else {
				panic(err)
			}
		} else {
			return createNonReadableArticle(storageDir, params, err)
		}
	} else {
		return createNonReadableArticle(storageDir, params, err)
	}
}

func createReadableArticle(storageDir string, params ArticleParams, article readability.Article) string {
	articleFilePath := ReadableArticleFilePath(storageDir, params, article)
	articleFile, _ := os.Create(articleFilePath)
	defer articleFile.Close()
	articleFile.WriteString(template.Render(template.ReadableArticleHTML(), template.CreateReadableArticleRenderArgs(params.Url, params.ArchivedAt, article)))
	return articleFilePath
}

func createNonReadableArticle(storageDir string, params ArticleParams, err error) string {
	articleFilePath := NonReadableArticleFilePath(storageDir, params)
	articleFile, _ := os.Create(articleFilePath)
	defer articleFile.Close()
	articleFile.WriteString(template.Render(template.NonReadableArticleHTML(), template.CreateNonReadableArticleRenderArgs(params.Url, params.ArchivedAt, err)))
	return articleFilePath
}

func createNonHTMLContent(storageDir string, params ArticleParams, contentType string, content []byte) string {
	contentFilePath := NonHTMLContentFilePath(storageDir, params, contentType)
	contentFile, _ := os.Create(contentFilePath)
	defer contentFile.Close()
	contentFile.Write(content)
	return contentFilePath
}
