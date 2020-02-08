package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	readability "github.com/go-shiori/go-readability"
)

const articlesRootDir = "articles"

func main() {
	url, err := ArticleURL(os.Args)
	if err != nil {
		//fmt.Println(err)
		fmt.Println("\nUSAGE: pocketize ARTICLE_URL")
		os.Exit(1)
	}

	os.MkdirAll(articlesRootDir, os.ModePerm)

	article, err := readability.FromURL(url.String(), 30*time.Second)
	if err == nil {
		fmt.Println(createArticle(url, article))
	} else {
		//fmt.Printf("failed to parse %s, %v\n", url, err)
		resp, err := http.Get(url.String())

		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()

			article, err := readability.FromReader(resp.Body, url.String())
			if err == nil {
				fmt.Println(createArticle(url, article))
			} else {
				//fmt.Printf("failed to parse %s: %v\n", url, err)
				fmt.Println(createArticleWithFailedReadability(url, err))
			}
		} else {
			//fmt.Printf("failed to download %s: %v\n", url, err)
			fmt.Println(createArticleWithFailedReadability(url, err))
		}
	}
}

func createArticle(url *url.URL, article readability.Article) string {
	articleFilePath := createArticleFilePath(url, article)
	articleFile, _ := os.Create(articleFilePath)
	defer articleFile.Close()
	args := RenderArgs{
		Address:     url,
		Title:       article.Title,
		Image:       article.Image,
		Excerpt:     article.Excerpt,
		Byline:      byline(article),
		SiteName:    siteName(url, article),
		ReadingTime: readingTime(article),
		Content:     template.HTML(article.Content),
		ArchivedAt:  time.Now().Format("January 2, 2006"),
	}

	articleFile.WriteString(Render(articleWithStyling(), args))
	return articleFilePath
}

func createArticleWithFailedReadability(url *url.URL, err error) string {
	articleFilePath := createArticleWithFailedReadabilityFilePath(url)
	articleFile, _ := os.Create(articleFilePath)
	defer articleFile.Close()
	args := RenderArgs{
		Address: url,
		Error:   err,
	}
	articleFile.WriteString(Render(articleWithFailedReadabilityWithStyling(), args))
	return articleFilePath
}

func createArticleWithFailedReadabilityFilePath(address *url.URL) string {
	return path.Join(articlesRootDir, fmt.Sprintf("%s-%s.html", time.Now().Format("2006-01-02"), formattedHost(address)))
}

func createArticleFilePath(address *url.URL, article readability.Article) string {
	return path.Join(articlesRootDir, fmt.Sprintf("%s-%s-%s.html", time.Now().Format("2006-01-02"), formattedTitle(article.Title), formattedHost(address)))
}

func formattedTitle(title string) string {
	replaceInvalidCharactersRegexp := regexp.MustCompile("[<>:\"'/\\|?*=;.%^ ]")
	replaceDuplicateAdjacentDashesRegexp := regexp.MustCompile("-{2,}")
	return replaceDuplicateAdjacentDashesRegexp.ReplaceAllString(replaceInvalidCharactersRegexp.ReplaceAllString(strings.ReplaceAll(title, "&", "and"), "-"), "-")
}

func formattedHost(address *url.URL) string {
	return strings.ReplaceAll(address.Host, ".", "-")
}

func articleWithStyling() string {
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

func articleWithFailedReadabilityWithStyling() string {
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
