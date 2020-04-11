package article

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func HttpEquivRefreshUrl(articleUrl *url.URL, doc string) *url.URL {
	isMetaRefreshTag := false
	var metaRefreshUrl *url.URL

	tokenizer := html.NewTokenizer(strings.NewReader(doc))
	for {
		if tokenizer.Next() == html.ErrorToken {
			break
		}

		if tagName, _ := tokenizer.TagName(); string(tagName) == "meta" {
			for {
				attrName, attrValue, hasMoreAttrs := tokenizer.TagAttr()
				if strings.ToLower(string(attrName)) == "http-equiv" && strings.ToLower(string(attrValue)) == "refresh" {
					isMetaRefreshTag = true
				} else if strings.ToLower(string(attrName)) == "content" {
					contentAttrValues := strings.Split(string(attrValue), ";")
					for _, value := range contentAttrValues {
						if strings.Contains(strings.ToLower(value), "url=") {
							possibleUrl := strings.ReplaceAll(strings.ReplaceAll(value, "url=", ""), "URL=", "")
							if refreshUrl, err := url.Parse(strings.TrimSpace(possibleUrl)); err == nil {
								metaRefreshUrl = refreshUrl
							}
						}
					}
				}

				if !hasMoreAttrs {
					break
				}
			}
		}

		if isMetaRefreshTag && metaRefreshUrl != nil {
			break
		}
	}

	if metaRefreshUrl != nil {
		if !metaRefreshUrl.IsAbs() {
			return articleUrl.ResolveReference(metaRefreshUrl)
		} else {
			return metaRefreshUrl
		}
	} else {
		return nil
	}
}
