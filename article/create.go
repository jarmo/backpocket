package article

import (
	"net/http"
	"net/url"
	"os"
	"time"
	"github.com/jarmo/pocketize/template"

	readability "github.com/go-shiori/go-readability"
)

func Create(url *url.URL) string {
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
