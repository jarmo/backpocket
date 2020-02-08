package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"

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

	article, err := readability.FromURL(url.String(), 30 * time.Second)
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

type renderArgs struct {
	Address     *url.URL
	Title       string
	Image       string
	Excerpt     string
	Byline      string
	SiteName    string
	ReadingTime int
	Content     template.HTML
	ArchivedAt  string
	Error       error
}

func createArticle(url *url.URL, article readability.Article) string {
	articleFilePath := createArticleFilePath(url, article)
	articleFile, _ := os.Create(articleFilePath)
	defer articleFile.Close()
	args := renderArgs{
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

	articleFile.WriteString(render(articleWithStyling(), args))
	return articleFilePath
}

func createArticleWithFailedReadability(url *url.URL, err error) string {
	articleFilePath := createArticleWithFailedReadabilityFilePath(url)
	articleFile, _ := os.Create(articleFilePath)
	defer articleFile.Close()
	args := renderArgs{
		Address: url,
		Error:   err,
	}
	articleFile.WriteString(render(articleWithFailedReadabilityWithStyling(), args))
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

func render(content string, args renderArgs) string {
	tmpl, err := template.New("article").Parse(content)
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
