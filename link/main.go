package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
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

	// Convert []byte to string, create new io.Reader, and Parse to Node Tree
	doc, err := html.Parse(strings.NewReader(string(h)))
	if err != nil {
		log.Fatalf("link: issue parsing html string from new reader: %v", err)
	}

	// Call function to parse Node Tree
	parseNodeTree(doc)
}

// Reference: https://pkg.go.dev/golang.org/x/net/html#Parse
func parseNodeTree(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				fmt.Println(a.Val)
				break
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseNodeTree(c)
	}
}
