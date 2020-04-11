package htmlparser

import (
	"bytes"
	"errors"
	"io"

	"golang.org/x/net/html"
)

func Render(node *html.Node) string {
	var buf bytes.Buffer
	html.Render(io.Writer(&buf), node)
	return buf.String()
}

type nodeForEachNodeFn func(*html.Node)

func ForEachNode(rootNode *html.Node, fn nodeForEachNodeFn) {
	walkNodes(rootNode, func(node *html.Node) error {
		fn(node)
		return nil
	})
}

type nodeFindFn func(*html.Node) bool

func FindNode(rootNode *html.Node, fn nodeFindFn) *html.Node {
	var foundNode *html.Node

	walkNodes(rootNode, func(node *html.Node) error {
		if fn(node) {
			foundNode = node
			return errors.New("interrupt")
		}

		return nil
	})

	return foundNode
}

func AttrByName(node *html.Node, name string) string {
	for _, attr := range node.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}

	return ""
}

type nodeWalkFn func(*html.Node) error

func walkNodes(rootNode *html.Node, fn nodeWalkFn) {
	var walker func(*html.Node)
	walker = func(node *html.Node) {
		err := fn(node)
		if err != nil {
			return
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walker(child)
		}
	}
	walker(rootNode)
}
