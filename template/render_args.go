package template

import (
	"html/template"
	"net/url"
	"strings"
	"time"

	readability "github.com/go-shiori/go-readability"
)

type RenderArgs struct {
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

func CreateReadableArticleRenderArgs(url *url.URL, article readability.Article) RenderArgs {
	return RenderArgs{
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
}

func CreateNonReadableArticleRenderArgs(url *url.URL, err error) RenderArgs {
	return RenderArgs{
		Address: url,
		Error:   err,
	}
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
