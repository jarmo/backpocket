package main

import (
	"net/http"
	"net/url"
	"os"
	"time"

	readability "github.com/go-shiori/go-readability"
)

func CreateArticle(url *url.URL) string {
	article, err := readability.FromURL(url.String(), 30*time.Second)
	if err == nil {
		return createReadableArticle(url, article)
	} else {
		//fmt.Printf("failed to parse %s, %v\n", url, err)
		resp, err := http.Get(url.String())

		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()

			article, err := readability.FromReader(resp.Body, url.String())
			if err == nil {
				return createReadableArticle(url, article)
			} else {
				//fmt.Printf("failed to parse %s: %v\n", url, err)
				return createNonReadableArticle(url, err)
			}
		} else {
			//fmt.Printf("failed to download %s: %v\n", url, err)
			return createNonReadableArticle(url, err)
		}
	}
}

func createReadableArticle(url *url.URL, article readability.Article) string {
	articleFilePath := ReadableArticleFilePath(url, article)
	articleFile, _ := os.Create(articleFilePath)
	defer articleFile.Close()
	articleFile.WriteString(Render(ReadableArticleHTML(), CreateReadableArticleRenderArgs(url, article)))
	return articleFilePath
}

func createNonReadableArticle(url *url.URL, err error) string {
	articleFilePath := NonReadableArticleFilePath(url)
	articleFile, _ := os.Create(articleFilePath)
	defer articleFile.Close()
	articleFile.WriteString(Render(NonReadableArticleHTML(), CreateNonReadableArticleRenderArgs(url, err)))
	return articleFilePath
}
