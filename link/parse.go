package link

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

var r io.Reader

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		log.Fatalf("link: issue parsing html string from new reader: %v", err)
	}

	parseNodeTree(doc)

	return nil, nil
}

// Reference: https://pkg.go.dev/golang.org/x/net/html#Parse
func parseNodeTree(n *html.Node) ([]Link, error) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				fmt.Println(a.Val, a.Namespace)
				break
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseNodeTree(c)
	}
	return nil, nil
}
