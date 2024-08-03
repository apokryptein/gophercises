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

func main() {
	site := flag.String("s", "", "site to crawl and map")
	outFile := flag.String("o", "sitemap.xml", "desired XML filename for output")
	depth := flag.Int("d", 1, "desired crawl depth")
	flag.Parse()

	if !isFlagPassed("s") {
		flag.Usage()
		os.Exit(1)
	}

	// Valide user input: site flag
	validateInput(site)

	// map the provided site
	sitemap := makeMapOfSite(*site, *depth)

	// create file for XML output -> io.Writer
	w, err := os.Create(*outFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sitemap: error creating file: %v", err)
		os.Exit(1)
	}

	// parse sitemap data into Url and UrlSet structs
	// for XML marshaling
	urls := makeUrlSlice(sitemap)
	uset := UrlSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Urls:  urls,
	}

	fmt.Printf("[i] Writing results to: %s\n", *outFile)

	// encode and write to file usign io.Writer
	w.WriteString(xml.Header)
	enc := xml.NewEncoder(w)
	enc.Indent("", "    ")
	if err := enc.Encode(uset); err != nil {
		fmt.Fprintf(os.Stderr, "sitemap: error encoding XML: %v", err)
		os.Exit(1)
	}
}

// Determine whether flag was passed
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// Validate user input contains Scheme -> http/s
func validateInput(s *string) {
	if strings.HasPrefix(*s, "https:") {
		return
	}
	fmt.Fprintf(os.Stderr, "sitemap: please supply full domain with scheme")
	os.Exit(1)
}

// Parse slice of string into slice of Url XML marshaling
func makeUrlSlice(sm []string) []Url {
	urls := make([]Url, 0)
	for _, l := range sm {
		urls = append(urls, Url{l})
	}
	return urls
}

// Map a given website at the specified depth
// Breadth-First Search (BFS) algorithm
func makeMapOfSite(seed string, depth int) []string {
	// tracks visited sites
	visited := make(map[string]struct{})

	// current and next queues
	var next map[string]struct{}
	queue := map[string]struct{}{
		seed: {},
	}

	for i := 0; i <= depth; i++ {
		next = make(map[string]struct{})
		for l := range queue {
			// if site in visited, skip and continue
			if _, ok := visited[l]; ok {
				continue
			}
			visited[l] = struct{}{}
			for _, link := range fetch(l) {
				// if link hasn't been visited add to next queue
				// otherwise continue
				if _, ok := visited[link]; !ok {
					next[link] = struct{}{}
				}
			}
		}
		// queue is next queue for next iteration
		queue = next
	}

	sitemap := make([]string, 0, len(visited))
	for l := range visited {
		sitemap = append(sitemap, l)
	}
	return sitemap
}

// Fetch links in a given page
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

	return parseUrls(resp)
}

// TODO: Possibly drop URLs containing '#'
// Fetch and parse links returning a []string of absolute URLs
func parseUrls(resp *http.Response) []string {
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
		switch {
		case strings.HasPrefix(l.Href, "/"):
			urls = append(urls, base+l.Href)
		case strings.HasPrefix(l.Href, "#"):
			urls = append(urls, base+"/"+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			urls = append(urls, l.Href)
		}
	}
	return filterScope(urls, base)
}

// Filters URLs to scope
func filterScope(links []string, base string) []string {
	var filtered []string
	for _, l := range links {
		if strings.HasPrefix(l, base) {
			filtered = append(filtered, l)
		}
	}
	return filtered
}
