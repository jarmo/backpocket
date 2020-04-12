package article

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	readability "github.com/go-shiori/go-readability"
)

func ReadableArticleFilePath(storageDir string, params ArticleParams, article readability.Article) string {
	return filepath.Join(storageDir, fmt.Sprintf("%s-%s-%s.html", params.ArchivedAt.Format("2006-01-02"), titleFromArticleOrPath(article, params.Url), suffix(params.Url)))
}

func NonReadableArticleFilePath(storageDir string, params ArticleParams) string {
	return filepath.Join(storageDir, fmt.Sprintf("%s-%s-%s.html", params.ArchivedAt.Format("2006-01-02"), titleFromUrl(params.Url), suffix(params.Url)))
}

func NonHTMLContentFilePath(storageDir string, params ArticleParams, contentType string) string {
	return filepath.Join(storageDir, fmt.Sprintf("%s-%s-%s.%s", params.ArchivedAt.Format("2006-01-02"), titleFromUrl(params.Url), suffix(params.Url), extension(contentType)))
}

func extension(contentType string) string {
	if strings.Contains(contentType, "text/plain") {
		return "txt"
	} else if strings.Contains(contentType, "xhtml") {
		return "html"
	} else {
		contentTypeParts := strings.Split(strings.Split(contentType, ";")[0], "/")
		if len(contentTypeParts) > 1 {
			return contentTypeParts[1]
		} else {
			return contentTypeParts[0]
		}
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

func suffix(address *url.URL) string {
	hash := sha1.New()
	hash.Write([]byte(address.String()))
	return hex.EncodeToString(hash.Sum(nil))[:8]
}
