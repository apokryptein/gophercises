package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/apokryptein/gophercises/link"
)

func main() {
	fmt.Println("Sitemap Utility")
	site := flag.String("site", "", "site to crawl and map")
	flag.Parse()

	if !isFlagPassed("site") {
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

	for _, link := range links {
		fmt.Println(link.Href, link.Text)
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
