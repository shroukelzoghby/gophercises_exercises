package main

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestExtractHref(t *testing.T) {
	cases := []struct {
		name string
		html string
		href string
	}{
		{
			name: "valid a element with href",
			html: `<a href="/example">Example</a>`,
			href: "/example",
		},
		{
			name: "missing href",
			html: `<a>Example</a>`,
			href: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			document, err := html.Parse(strings.NewReader(c.html))
			if err != nil {
				t.Fatalf("Test %s failed: error parsing HTML: %v", c.name, err)
			}
			anchors := collectAnchors(document)
			if len(anchors) != 1 {
				t.Fatalf("expected one anchor, got %d", len(anchors))
			}
			href := extractHref(anchors[0])
			if href != c.href {
				t.Fatalf("Test %s failed: expected href %s, got %s", c.name, c.href, href)
			}
		})

	}
}

func TestExtractText(t *testing.T) {
	document, err := html.Parse(strings.NewReader(`
		<a href="/dog">
			<span>Something in a span</span>
			Text not in a span
			<!-- This comment should be ignored. -->
			<b>Bold text!</b>
		</a>
	`))
	if err != nil {
		t.Fatalf("error parsing HTML: %v", err)
	}

	anchors := collectAnchors(document)
	if got := extractText(anchors[0]); got != "Something in a span Text not in a span Bold text!" {
		t.Fatalf("expected normalized link text, got %q", got)
	}
}

func TestFindAnchorsIgnoresNestedLinks(t *testing.T) {
	inner := &html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{{Key: "href", Val: "/inner"}},
	}
	outer := &html.Node{
		Type:       html.ElementNode,
		Data:       "a",
		Attr:       []html.Attribute{{Key: "href", Val: "/outer"}},
		FirstChild: inner,
	}
	inner.Parent = outer
	document := &html.Node{Type: html.DocumentNode, FirstChild: outer}
	outer.Parent = document

	anchors := collectAnchors(document)
	if len(anchors) != 1 {
		t.Fatalf("expected one outer anchor, got %d", len(anchors))
	}
	if got := extractHref(anchors[0]); got != "/outer" {
		t.Fatalf("expected outer href /outer, got %q", got)
	}
}

func collectAnchors(document *html.Node) []*html.Node {
	anchors := make(chan *html.Node)
	go findAnchors(document, anchors)

	var result []*html.Node
	for anchor := range anchors {
		result = append(result, anchor)
	}
	return result
}
