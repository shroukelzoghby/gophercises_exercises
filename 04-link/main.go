package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func main() {

	flagHtml := flag.String("html", "links.html", "Path to HTML file containing the links")
	flag.Parse()
	file, err := os.Open(*flagHtml)
	if err != nil {
		log.Fatalf("Error opening HTML file: %v", err)
	}
	root, err := html.Parse(file)
	if err != nil {
		log.Fatalf("Error parsing HTML file: %v", err)
	}

	aChan := make(chan *html.Node)
	go findAnchors(root, aChan)
	for a := range aChan {
		fmt.Println(Link{
			Href: extractHref(a),
			Text: extractText(a),
		})
	}

}

func findAnchors(n *html.Node, aChan chan *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		aChan <- n
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findAnchors(c, aChan)
	}

	if n.Parent == nil {
		close(aChan)
	}
}

func extractHref(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}
func extractText(n *html.Node) string {
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			text += c.Data
			continue
		}
		text += extractText(c)
	}
	return strings.Join(strings.Fields(text), " ")
}
