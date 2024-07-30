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
	validateSite(site)

	links := fetchLinks(*site)

	sitemap := makeMapOfSite(links, *site, *depth)
	printSiteMap(sitemap)

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

// Test user input to ensure Scheme is present
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

// Parse slice of Link into slice of Url
// for XML marshaling
func makeUrlSlice(sm map[int][]link.Link) []Url {
	// origDomain := getDomain(s)

	urls := make([]Url, 0, getMapSize((sm)))

	for _, l := range sm {
		for _, u := range l {
			urls = append(urls, Url{u.Href})
		}
	}
	return urls
}

// Function to map a given website at the specified depth
func makeMapOfSite(seed []link.Link, site string, depth int) map[int][]link.Link {
	sitemap := make(map[int][]link.Link)
	visited := make(map[string]struct{})
	// TODO: add queue, visited not currently fully functional for BFS

	for i := range depth {
		for _, link := range seed {
			url := makeUrl(link.Href, site)
			if _, ok := visited[url]; ok {
				continue
			}
			links := fetchLinks(url)
			updatedLinks := updateUrls(links, site)
			sitemap[i] = append(sitemap[i], updatedLinks...)
			visited[url] = struct{}{}
		}
	}
	return sitemap
}

// Returns path of given URL
func getPath(site string) string {
	url, err := url.Parse(site)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sitemap: error parsing URL: %v", err)
	}

	path := url.Path
	return path
}

// Return domain name of given URL
// For example: https://www.google.com returns -> "google"
// Used to test domain name to remain in crawl scope
func getDomain(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sitemap: error parsing URL: %v", err)
		os.Exit(1)
	}

	url := strings.Split(u.Hostname(), ".")
	domain := url[len(url)-2]

	return domain
}

// Function to fetch links in a given page
func fetchLinks(s string) []link.Link {
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

	h := resp.Body

	links, err := link.Parse(h)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sitemap: error parsing links: %v", err)
		os.Exit(1)
	}
	return links
}

// Function to return an Absolute URL for a given input
func makeUrl(currentLink string, domain string) string {
	url, err := url.Parse(currentLink)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sitemap: error parsing url: %v", err)
		os.Exit(1)
	}

	if url.Scheme == "" {
		d, _ := url.Parse(domain)
		if strings.HasPrefix(currentLink, "#") {
			return "https://" + d.Host + "/" + currentLink
		}
		return "https://" + d.Host + currentLink
	}

	return currentLink
}

// Function to update all URLs in a given []link.Link
func updateUrls(links []link.Link, domain string) []link.Link {
	updatedLinks := make([]link.Link, 0, len(links))
	for _, l := range links {
		ul := makeUrl(l.Href, domain)
		if !isLinkInScope(ul, domain) {
			continue
		}
		updatedLinks = append(updatedLinks, link.Link{Href: ul, Text: l.Text})
	}
	return updatedLinks
}

// Helper function to print sitemap results
func printSiteMap(sm map[int][]link.Link) {
	fmt.Printf("DEPTH: %d\n", len(sm))
	for _, i := range sm {
		for _, link := range i {
			fmt.Println(link.Href)
		}
	}
}

// Determines whether link is in scope
func isLinkInScope(l string, s string) bool {
	ls := getDomain(l)
	ss := getDomain(s)

	return ls == ss
}

// Returns number of elements in map[int][]link.Link
func getMapSize(m map[int][]link.Link) int {
	size := 0

	for _, v := range m {
		size += len(v)
	}

	return size
}
