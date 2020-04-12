package article

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"math/rand"
	"math"

	readability "github.com/go-shiori/go-readability"
)

func ReadableArticleFilePath(storageDir string, params ArticleParams, article readability.Article) string {
	return filepath.Join(storageDir, fmt.Sprintf("%s-%s-%s.html", params.ArchivedAt.Format("2006-01-02"), titleFromArticleOrPath(article, params.Url), randomSuffix()))
}

func NonReadableArticleFilePath(storageDir string, params ArticleParams) string {
	return filepath.Join(storageDir, fmt.Sprintf("%s-%s-%s.html", params.ArchivedAt.Format("2006-01-02"), titleFromUrl(params.Url), randomSuffix()))
}

func NonHTMLContentFilePath(storageDir string, params ArticleParams, contentType string) string {
	return filepath.Join(storageDir, fmt.Sprintf("%s-%s-%s.%s", params.ArchivedAt.Format("2006-01-02"), titleFromUrl(params.Url), randomSuffix(), extension(contentType)))
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
		return titleFromUrl(url)
	}
}

func titleFromUrl(url *url.URL) string {
	urlPath := url.Path

	if len(urlPath) > 1 {
		return formattedTitle(strings.ReplaceAll(filepath.Base(urlPath), filepath.Ext(urlPath), ""))
	} else {
		return formattedHost(url)
	}
}

func formattedTitle(title string) string {
	replaceInvalidCharactersRegexp := regexp.MustCompile("[^\x00-\x7F]")
	replaceUnsupportedCharactersRegexp := regexp.MustCompile("[<>:\"'/\\|?*=;.%,^{}]")
	replaceDuplicateAdjacentDashesRegexp := regexp.MustCompile("-{2,}")
	formattedTitle := strings.TrimSpace(strings.Trim(strings.Trim(
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
	
	maxTitleLength := int(math.Min(64, float64(len(formattedTitle))))
	return formattedTitle[:maxTitleLength]
}

func formattedHost(address *url.URL) string {
	return strings.ReplaceAll(address.Hostname(), ".", "-")
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
