package main

import (
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/shroukelzoghby/gophercises_exercises/sitemap/link"
)

func main() {
	flagURL := flag.String("url", "", "The URL to create a sitemap for.")
	flagDepth := flag.Int("depth", 2, "The maximum depth of links to follow.")
	flagXMLFilename := flag.String("xml", "sitemap.xml", "The output sitemap XML filename.")
	flag.Parse()

	if *flagURL == "" {
		log.Fatal("Error: Missing -url flag. Usage: sitemap -url https://example.com [-depth 2] [-xml sitemap.xml]")
	}

	if *flagDepth <= 0 {
		log.Fatal("Error: -depth must be greater than 0")
	}

	sitemap, err := buildSitemap(*flagURL, *flagDepth)
	if err != nil {
		log.Fatalf("Error: Failed to build sitemap for %s: %v", *flagURL, err)
	}

	if len(sitemap) == 0 {
		log.Fatalf("Error: No URLs found for %s", *flagURL)
	}

	if err := generateSitemap(sitemap, *flagXMLFilename); err != nil {
		log.Fatalf("Error: Failed to write sitemap to %s: %v", *flagXMLFilename, err)
	}

	log.Printf("Successfully generated sitemap with %d URL(s) for %s in %s", len(sitemap), *flagURL, *flagXMLFilename)
}

func buildSitemap(baseURL string, depth int) ([]string, error) {
	visited := make(map[string]bool)
	allURLs := []string{}
	currentLevel := []string{baseURL}

	for d := 0; d < depth; d++ {
		var nextLevel []string

		for _, url := range currentLevel {
			if visited[url] {
				continue
			}
			visited[url] = true
			allURLs = append(allURLs, url)

			subURLs, err := getURLs(url)
			if err != nil {
				log.Printf("Warning: Failed to fetch URLs from %s: %v", url, err)
				continue
			}

			for _, subURL := range subURLs {
				if !visited[subURL] {
					nextLevel = append(nextLevel, subURL)
				}
			}
		}

		if len(nextLevel) == 0 {
			break
		}
		currentLevel = nextLevel
	}

	return allURLs, nil
}

func getURLs(pageURL string) ([]string, error) {
	pageURL = strings.TrimSuffix(pageURL, "/")

	res, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	links, err := link.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	var domainURLs []string
	for _, l := range links {
		normalizedURL := normalizeURL(l.Href, pageURL)
		if normalizedURL != "" {
			domainURLs = append(domainURLs, normalizedURL)
		}
	}

	return domainURLs, nil
}

// normalizeURL converts a link href to an absolute URL if it belongs to the same domain.
// Returns empty string if the URL should be filtered out.
func normalizeURL(href, baseURL string) string {
	// Skip external URLs
	if strings.HasPrefix(href, "http") && !strings.HasPrefix(href, baseURL) {
		return ""
	}

	// Already a full domain URL
	if strings.HasPrefix(href, baseURL) {
		return href
	}

	// Skip protocol-relative and other non-path URLs
	if strings.Contains(href, ":") {
		return ""
	}

	// Remove fragment identifiers
	if idx := strings.Index(href, "#"); idx != -1 {
		href = href[:idx]
	}

	// Handle empty URLs
	if href == "" {
		return ""
	}

	// Ensure path starts with /
	if href[0] != '/' {
		href = "/" + href
	}

	return baseURL + href
}

type SitemapXML struct {
	XMLName xml.Name        `xml:"urlset"`
	Xmlns   string          `xml:"xmlns,attr"`
	URLs    []SitemapXMLURL `xml:"url"`
}

type SitemapXMLURL struct {
	Loc string `xml:"loc"`
}

func generateSitemap(urls []string, filename string) error {
	sitemap := buildSitemapXML(urls)

	xmlData, err := marshalSitemap(sitemap)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, xmlData, 0644)
}

func buildSitemapXML(urls []string) SitemapXML {
	sitemap := SitemapXML{
		Xmlns: "https://www.sitemaps.org/schemas/sitemap/0.9",
	}
	for _, url := range urls {
		sitemap.URLs = append(sitemap.URLs, SitemapXMLURL{Loc: url})
	}
	return sitemap
}

func marshalSitemap(sitemap SitemapXML) ([]byte, error) {
	xmlData, err := xml.MarshalIndent(&sitemap, "", "\t")
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), xmlData...), nil
}
