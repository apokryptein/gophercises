package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"os"

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

	urls := make([]Url, len(links))

	for i, link := range links {
		urls[i] = Url{Loc: link.Href}
	}

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
