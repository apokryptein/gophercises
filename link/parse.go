package link

import (
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

// var r io.Reader
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		log.Fatalf("link: issue parsing html string from new reader: %v", err)
	}

	linkNodes := parseNodeTree(doc)

	var links []Link
	for _, n := range linkNodes {
		links = append(links, createLink(n))
	}
	return links, nil
}

// Reference: https://pkg.go.dev/golang.org/x/net/html#Parse
// Implements Depth First Search (DFS) algorithm
func parseNodeTree(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		// fmt.Println(a.Val, a.Namespace)
		return []*html.Node{n}
	}

	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, parseNodeTree(c)...)
	}
	return ret
}

// creates a Link struct from a given node
func createLink(n *html.Node) Link {
	var l Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			l.Href = attr.Val
			break
		}
	}
	l.Text = formatText(retrieveText(n))
	return l
}

func retrieveText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	if n.Type != html.ElementNode {
		return ""
	}

	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += retrieveText(c) + " "
	}
	return ret
}

func formatText(s string) string {
	sl := strings.Fields(s)
	return strings.Join(sl, " ")
}
