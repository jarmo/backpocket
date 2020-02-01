package main

import (
	"fmt"
	"os"
	"time"
	"strings"
	"path"
	"regexp"
	"io/ioutil"
	"net/http"
	"net/url"
	"encoding/base64"
	"golang.org/x/net/html"

	readability "github.com/go-shiori/go-readability"
)

func main() {
	uri, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Printf("failed to parse url %v\n", err)
		os.Exit(1)
	}

	replaceSchemeRegexp := regexp.MustCompile("http(s)?://")
	replaceInvalidCharactersRegexp := regexp.MustCompile("[<>:\"\\|?*=&;]")
	articleDir := "articles/" + replaceInvalidCharactersRegexp.ReplaceAllString(replaceSchemeRegexp.ReplaceAllString(uri.String(), ""), "_")
	os.MkdirAll(articleDir, os.ModePerm)

	articlePath := path.Join(articleDir, "index.html")

	articleFile, _ := os.Create(articlePath)
	defer articleFile.Close()

	article, err := readability.FromURL(uri.String(), 30 * time.Second)
	if err != nil {
		//fmt.Printf("failed to parse %s, %v\n", uri, err)
		resp, err := http.Get(uri.String())
		defer resp.Body.Close()

		if err != nil {
			//fmt.Printf("failed to download %s: %v\n", uri, err)
			articleFile.WriteString(addressFile(uri, err))
		} else {
			article, err := readability.FromReader(resp.Body, uri.String())
			if err != nil {
				//fmt.Printf("failed to parse %s: %v\n", uri, err)
				articleFile.WriteString(addressFile(uri, err))
			} else {
				articleFile.WriteString(articleWithStyling(uri, article))
			}
		}
	} else {
		articleFile.WriteString(articleWithStyling(uri, article))
	}

	fmt.Println(articlePath)
}

func articleWithStyling(uri *url.URL, article readability.Article) string {
	archivedAt := time.Now().UTC()

	return contentWithBase64Images(fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
			<head>
				<meta content="text/html;charset=utf-8" http-equiv="Content-Type">
				<meta content="utf-8" http-equiv="encoding">
				<title>%s</title>
				<style>%s</style>
			</head>
			<body>
				<header>
					<h1>
						<a href="%s">%s</a>
						<div class="archived-at">%s</div>
					</h1>
					<img src="%s">
					<figcaption>%s</figcaption>
					<small>%s • %s • %d minutes</small>
				</header>
				<article data-archived-at="%s">%s</article>
			</body>
		</html>
		`,
		article.Title,
		Styles(),
		uri.String(), article.Title,
		archivedAt.Format("January 2, 2006"),
		article.Image,
		article.Excerpt,
		byline(article), siteName(uri, article), readingTime(article),
		archivedAt.Format(time.RFC3339), article.Content))
}

func addressFile(address *url.URL, err error) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
	<head>
		<meta content="text/html;charset=utf-8" http-equiv="Content-Type">
		<meta content="utf-8" http-equiv="encoding">
		<title>%s</title>
		<style>%s</style>
	</head>
	<body>
	  <header>
			<h1><a href="%s">%s</a></h1>
			<figcaption>%v</figcaption>
	  </header>
	</body>
</html>
	`, address.String(), Styles(), address.String(), address.String(), err)
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

func contentWithBase64Images(doc string) string {
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

