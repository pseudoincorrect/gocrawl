package parser

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

func Parse(htmlDoc io.Reader) ([]string, error) {
	doc, err := html.Parse(htmlDoc)
	if err != nil {
		return nil, fmt.Errorf("html parsing err : %v", err)
	}
	var links []string = nil
	visit(links, doc)
	return links, nil
}

func visit(links []string, n *html.Node) {
	// printNode(n)
	if n.Type == html.ElementNode && n.Data == "a" {
		link := n.Attr[0].Val
		links = append(links, link)
		fmt.Println(link)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		visit(links, c)
	}
}

func printNode(n *html.Node) {
	fmt.Println("n.Type = ", n.Type)
	fmt.Println("n.Data = ", n.Data)
	for _, a := range n.Attr {
		fmt.Println("attr = ", a)
	}
}
