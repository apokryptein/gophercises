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

// TODO: Add functionality to skip already seen links in sitemap data structure
// TODO: REFACTOR
func main() {
	site := flag.String("s", "", "site to crawl and map")
	outFile := flag.String("o", "sitemap.xml", "desired XML filename for output")
	depth := flag.Int("d", 1, "desired crawl depth")
	flag.Parse()

	if !isFlagPassed("s") {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("You have chosen to map: ", *site)

	// validate user input contains Scheme -> http/s
	validateInput(site)

	// TODO: address case where depth is 1
	// currently minimum is set to 2 due to fetch prior
	// to pass to makeMapOfSite
	links := fetch(*site)

	sitemap := makeMapOfSite(links, *depth)

	w, err := os.Create(*outFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sitemap: error creating file: %v", err)
		os.Exit(1)
	}

	urls := makeUrlSlice(sitemap)

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

// Function to determine whether flag was passed
// on command line
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// TODO: get rid of this?
// Test user input to ensure Scheme is present
func validateInput(s *string) {
	if strings.HasPrefix(*s, "https:") {
		return
	}
	fmt.Fprintf(os.Stderr, "sitemap: please supply full domain with scheme")
	os.Exit(1)
}

// Parse slice of Link into slice of Url
// for XML marshaling
func makeUrlSlice(sm map[int][]string) []Url {
	urls := make([]Url, 0, getMapSize((sm)))

	for _, l := range sm {
		for _, u := range l {
			urls = append(urls, Url{u})
		}
	}
	return urls
}

// Function to map a given website at the specified depth
func makeMapOfSite(seed []string, depth int) map[int][]string {
	sitemap := make(map[int][]string)
	visited := make(map[string]struct{})
	// TODO: add queue, visited not currently fully functional for BFS

	for i := range depth {
		for _, l := range seed {
			// url := makeUrl(link.Href)
			if _, ok := visited[l]; ok {
				continue
			}
			links := fetch(l)
			// updatedLinks := updateUrls(links, site)
			// sitemap[i] = append(sitemap[i], updatedLinks...)
			sitemap[i] = append(sitemap[i], links...)
			visited[l] = struct{}{}
		}
	}
	return sitemap
}

// Function to fetch links in a given page
func fetch(s string) []string {
	resp, err := http.Get(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sitemap: error making get request: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "sitemap: status not 200")
		os.Exit(1)
	}

	return makeUrl(resp)
}

func makeUrl(resp *http.Response) []string {
	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()

	r := resp.Body

	links, err := link.Parse(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sitemap: error parsing links: %v", err)
		os.Exit(1)
	}

	var urls []string
	for _, l := range links {
		if strings.HasPrefix(l.Href, "/") {
			urls = append(urls, base+l.Href)
		} else if strings.HasPrefix(l.Href, "#") {
			urls = append(urls, base+"/"+l.Href)
		} else if strings.HasPrefix(l.Href, "http") {
			urls = append(urls, l.Href)
		}
	}
	return filterScope(urls, base)
}

// filters URLs to scope to site
func filterScope(links []string, base string) []string {
	var filtered []string
	for _, l := range links {
		if strings.HasPrefix(l, base) {
			filtered = append(filtered, l)
		}
	}
	return filtered
}

// Returns number of elements in map[int][]link.Link
func getMapSize(m map[int][]string) int {
	size := 0

	for _, v := range m {
		size += len(v)
	}

	return size
}
