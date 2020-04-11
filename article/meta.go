package article

import (
	"net/url"
	"strings"

	"github.com/jarmo/backpocket/htmlparser"

	"golang.org/x/net/html"
)

func HttpEquivRefreshUrl(articleUrl *url.URL, content string) *url.URL {
	doc, _ := html.Parse(strings.NewReader(content))

	metaRefreshNode := htmlparser.FindNode(doc, func(node *html.Node) bool {
		return node.Type == html.ElementNode && node.Data == "meta" && htmlparser.AttrByName(node, "http-equiv") != "" && htmlparser.AttrByName(node, "content") != ""
	})

	if metaRefreshNode != nil {
		contentAttrValue := htmlparser.AttrByName(metaRefreshNode, "content")
		contentAttrValueParts := strings.Split(contentAttrValue, ";")
		for _, value := range contentAttrValueParts {
			if strings.Contains(strings.ToLower(value), "url=") {
				possibleUrl := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(value, "url=", ""), "URL=", ""))
				if refreshUrl, err := url.Parse(possibleUrl); err == nil {
					return refreshUrl
				}
			}
		}
	}

	return nil
}

