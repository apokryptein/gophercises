package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/apokryptein/gophercises/link"
)

type Url struct {
	Loc string `xml:"loc"`
}

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []Url    `xml:"url"`
}

// TODO: now that it's functional, refactor this mess
func main() {
	site := flag.String("s", "", "site to crawl and map")
	outFile := flag.String("o", "sitemap.xml", "desired XML filename for output")
	flag.Parse()

	if !isFlagPassed("s") {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("You have chosen to map: ", *site)

	validateSite(site)
	resp, err := http.Get(*site)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sitemap: error making get request: %v", err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "sitemap: status not 200")
		os.Exit(1)
	}

	h := resp.Body

	links, _ := link.Parse(h)

	w, err := os.Create(*outFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sitemap: error creating file: %v", err)
		os.Exit(1)
	}

	urls := makeUrlSlice(links, *site)

	uset := UrlSet{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9", Urls: urls}

	fmt.Printf("[i] Writing results to: %s\n", *outFile)
	w.WriteString(xml.Header)
	enc := xml.NewEncoder(w)
	enc.Indent("", "    ")
	if err := enc.Encode(uset); err != nil {
		fmt.Fprintf(os.Stderr, "sitemap: error encoding XML: %v", err)
		os.Exit(1)
	}
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func validateSite(s *string) {
	if strings.HasPrefix(*s, "https://") {
		return
	} else if strings.HasPrefix(*s, "http://") {
		n := strings.Replace(*s, "http:", "https:", 1)
		fmt.Println("[!] HTTP not supported. Updating to HTTPS.")
		*s = n
		return
	}

	*s = "https://" + *s
	fmt.Printf("[i] URL updated: %s/n", *s)
}

// TODO: make this more modular
func makeUrlSlice(l []link.Link, s string) []Url {
	ou, err := url.Parse(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sitemap: error parsing URL: %v", err)
		os.Exit(1)
	}
	parts := strings.Split(ou.Hostname(), ".")
	origDomain := parts[len(parts)-2]

	urls := make([]Url, 0, len(l))

	for _, link := range l {
		cu, err := url.Parse(link.Href)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sitemap: error parsing URL: %v", err)
			os.Exit(1)
		}

		if cu.Scheme != "" {
			cparts := strings.Split(cu.Hostname(), ".")
			currDomain := cparts[len(cparts)-2]
			if currDomain != origDomain {
				continue
			}
		}
		if link.Href != "" {
			urls = append(urls, Url{Loc: link.Href})
		}
	}
	return urls
}
