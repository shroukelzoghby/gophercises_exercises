package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Text string
	Href string
}

func Parse(r io.Reader) ([]Link, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return extractLinks(root), nil
}

func extractLinks(n *html.Node) []Link {
	var links []Link
	traverseAnchorNodes(n, func(a *html.Node) {
		links = append(links, Link{
			Text: extractText(a),
			Href: extractHref(a),
		})
	})
	return links
}

func traverseAnchorNodes(n *html.Node, callback func(*html.Node)) {
	if n.Type == html.ElementNode && n.Data == "a" {
		callback(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverseAnchorNodes(c, callback)
	}
}

func extractText(a *html.Node) string {
	var text string
	for c := a.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			text += c.Data
			continue
		}
		text += extractText(c)
	}
	return strings.TrimSpace(text)
}

func extractHref(a *html.Node) string {
	for _, attr := range a.Attr {
		if attr.Key != "href" {
			continue
		}
		return attr.Val
	}
	return ""
}
