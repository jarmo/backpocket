package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	articleFile.WriteString(Render(readableArticleHTML(), CreateReadableArticleRenderArgs(url, article)))
	return articleFilePath
}

func createNonReadableArticle(url *url.URL, err error) string {
	articleFilePath := NonReadableArticleFilePath(url)
	articleFile, _ := os.Create(articleFilePath)
	defer articleFile.Close()
	articleFile.WriteString(Render(nonReadableArticleHTML(), CreateNonReadableArticleRenderArgs(url, err)))
	return articleFilePath
}

func readableArticleHTML() string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
			<head>
				<meta content="text/html;charset=utf-8" http-equiv="Content-Type">
				<meta content="utf-8" http-equiv="encoding">
				<title>{{.Title}}</title>
				<style>%s</style>
			</head>
			<body>
				<header>
					<h1>
						<a href="{{.Address}}">{{.Title}}</a>
						<div class="archived-at">Archived at {{.ArchivedAt}}</div>
					</h1>
					<img src="{{.Image}}">
					<figcaption>{{.Excerpt}}</figcaption>
					<small>{{.Byline}} • {{.SiteName}} • {{.ReadingTime}} minutes</small>
				</header>
				<article>{{.Content}}</article>
			</body>
		</html>
		`, Styles())
}

func nonReadableArticleHTML() string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
	<head>
		<meta content="text/html;charset=utf-8" http-equiv="Content-Type">
		<meta content="utf-8" http-equiv="encoding">
		<title>{{.Address}}</title>
		<style>%s</style>
	</head>
	<body>
	  <header>
			<h1><a href="{{.Address}}">{{.Address}}</a></h1>
			<figcaption>{{.Error}}</figcaption>
	  </header>
	</body>
</html>
	`, Styles())
}

func byline(article readability.Article) string {
	if len(article.Byline) > 0 {
		return article.Byline
	} else {
		return "N/A"
	}
}

func siteName(address *url.URL, article readability.Article) string {
	if len(article.SiteName) > 0 {
		return article.SiteName
	} else {
		return address.Host
	}
}

func readingTime(article readability.Article) int {
	wordsPerMinuteAverageReadingRate := 200
	return len(strings.Split(article.TextContent, " ")) / wordsPerMinuteAverageReadingRate
}
