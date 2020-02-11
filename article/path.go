package article

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	readability "github.com/go-shiori/go-readability"
)

const RootDir = "backpocket-contents"

func ReadableArticleFilePath(params ArticleParams, article readability.Article) string {
	return path.Join(RootDir, fmt.Sprintf("%s-%s.html", params.ArchivedAt.Format("2006-01-02"), formattedTitle(article.Title)))
}

func NonReadableArticleFilePath(params ArticleParams) string {
	return path.Join(RootDir, fmt.Sprintf("%s-%s.html", params.ArchivedAt.Format("2006-01-02"), titleFromPath(params.Url)))
}

func NonHTMLContentFilePath(params ArticleParams, contentType string) string {
	return path.Join(RootDir, fmt.Sprintf("%s-%s.%s", params.ArchivedAt.Format("2006-01-02"), titleFromPath(params.Url), extension(contentType)))
}

func extension(contentType string) string {
	if strings.Contains(contentType, "text/plain") {
		return "txt"
	} else {
		return strings.Split(strings.Split(contentType, ";")[0], "/")[1]
	}
}

func titleFromPath(url *url.URL) string {
	return formattedTitle(strings.ReplaceAll(path.Base(url.Path), path.Ext(url.Path), ""))
}

func formattedTitle(title string) string {
	replaceInvalidCharactersRegexp := regexp.MustCompile("[<>:\"'/\\|?*=;.%,^]")
	replaceDuplicateAdjacentDashesRegexp := regexp.MustCompile("-{2,}")
	return replaceDuplicateAdjacentDashesRegexp.ReplaceAllString(replaceInvalidCharactersRegexp.ReplaceAllString(strings.ReplaceAll(strings.ReplaceAll(title, " ", "-"), "&", "and"), ""), "-")
}

func formattedHost(address *url.URL) string {
	return strings.ReplaceAll(address.Host, ".", "-")
}
