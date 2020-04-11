package article

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"
	"math/rand"

	readability "github.com/go-shiori/go-readability"
)

const RootDir = "backpocket-contents"

func ReadableArticleFilePath(params ArticleParams, article readability.Article) string {
	return path.Join(RootDir, fmt.Sprintf("%s-%s-%s.html", params.ArchivedAt.Format("2006-01-02"), titleFromArticleOrPath(article, params.Url), randomSuffix()))
}

func NonReadableArticleFilePath(params ArticleParams) string {
	return path.Join(RootDir, fmt.Sprintf("%s-%s-%s.html", params.ArchivedAt.Format("2006-01-02"), titleFromPath(params.Url), randomSuffix()))
}

func NonHTMLContentFilePath(params ArticleParams, contentType string) string {
	return path.Join(RootDir, fmt.Sprintf("%s-%s-%s.%s", params.ArchivedAt.Format("2006-01-02"), titleFromPath(params.Url), randomSuffix(), extension(contentType)))
}

func extension(contentType string) string {
	if strings.Contains(contentType, "text/plain") {
		return "txt"
	} else {
		return strings.Split(strings.Split(contentType, ";")[0], "/")[1]
	}
}

func titleFromArticleOrPath(article readability.Article, url *url.URL) string {
	if len(strings.TrimSpace(article.Title)) > 0 {
		return formattedTitle(article.Title)
	} else {
		return titleFromPath(url)
	}
}

func titleFromPath(url *url.URL) string {
	return formattedTitle(strings.ReplaceAll(path.Base(url.Path), path.Ext(url.Path), ""))
}

func formattedTitle(title string) string {
	replaceInvalidCharactersRegexp := regexp.MustCompile("[^\x00-\x7F]")
	replaceUnsupportedCharactersRegexp := regexp.MustCompile("[<>:\"'/\\|?*=;.%,^{}]")
	replaceDuplicateAdjacentDashesRegexp := regexp.MustCompile("-{2,}")
	return strings.TrimSpace(strings.Trim(strings.Trim(
		replaceDuplicateAdjacentDashesRegexp.ReplaceAllString(
			replaceUnsupportedCharactersRegexp.ReplaceAllString(
				replaceInvalidCharactersRegexp.ReplaceAllString(
					strings.ReplaceAll(
						strings.ReplaceAll(title,
						" ", "-"),
					"&", "and"),
				""),
			""),
		"-"),
	"-"),
	"."))
}

func formattedHost(address *url.URL) string {
	return strings.ReplaceAll(address.Host, ".", "-")
}

func randomSuffix() string {
	rand.Seed(time.Now().UnixNano())
	const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
			b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}
