package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/apokryptein/gophercises/link"
)

func main() {
	// Obtain HTML file to parse from command line
	filename := flag.String("f", "", "file containing HTML")
	flag.Parse()

	// Ensure a filename was provided
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	// Read HTML file in as []byte
	h, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatalf("link: problem reading html file: %v", err)
		os.Exit(1)
	}

	// Call function to parse Node Tree
	links, _ := link.Parse(strings.NewReader(string(h)))

	for _, link := range links {
		fmt.Printf("HREF: %s\tTEXT: %s\n", link.Href, link.Text)
	}
}
