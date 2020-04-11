package htmlparser

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

func Render(node *html.Node) string {
	var buf bytes.Buffer
	html.Render(io.Writer(&buf), node)
	return buf.String()
}

type nodeFn func(*html.Node)

func ForEachNode(rootNode *html.Node, fn nodeFn) {
	var walker func(*html.Node)
	walker = func(node *html.Node) {
		fn(node)
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walker(child)
		}
	}
	walker(rootNode)
}

func AttrByName(node *html.Node, name string) string {
	for _, attr := range node.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}

	return ""
}
