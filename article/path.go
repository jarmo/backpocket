package article

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	readability "github.com/go-shiori/go-readability"
)

const RootDir = "articles"

func ReadableArticleFilePath(address *url.URL, article readability.Article) string {
	return path.Join(RootDir, fmt.Sprintf("%s-%s-%s.html", time.Now().Format("2006-01-02"), formattedTitle(article.Title), formattedHost(address)))
}

func NonReadableArticleFilePath(address *url.URL) string {
	return path.Join(RootDir, fmt.Sprintf("%s-%s.html", time.Now().Format("2006-01-02"), formattedHost(address)))
}

func NonHTMLContentFilePath(address *url.URL, contentType string) string {
	return path.Join(RootDir, fmt.Sprintf("%s-%s.%s", time.Now().Format("2006-01-02"), formattedTitle(strings.ReplaceAll(path.Base(address.Path), path.Ext(address.Path), "")), extension(contentType)))
}

func extension(contentType string) string {
	if strings.Contains(contentType, "text/plain") {
		return "txt"
	} else {
		return strings.Split(strings.Split(contentType, ";")[0], "/")[1]
	}
}

func formattedTitle(title string) string {
	replaceInvalidCharactersRegexp := regexp.MustCompile("[<>:\"'/\\|?*=;.%^ ]")
	replaceDuplicateAdjacentDashesRegexp := regexp.MustCompile("-{2,}")
	return replaceDuplicateAdjacentDashesRegexp.ReplaceAllString(replaceInvalidCharactersRegexp.ReplaceAllString(strings.ReplaceAll(title, "&", "and"), "-"), "-")
}

func formattedHost(address *url.URL) string {
	return strings.ReplaceAll(address.Host, ".", "-")
}
