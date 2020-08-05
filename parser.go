package linkparser

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link will parse the HTML <a> tag
// and converts it into a valid form
type Link struct {
	ID   int
	Href string
	Text string
}

// Parse the given html and return links
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	// fmt.Printf("%s", doc.Data)
	if err != nil {
		panic(err)
	}
	linkNodes := linkParse(doc)
	fmt.Println("Found ", len(linkNodes), " links")
	var links []Link
	for i, node := range linkNodes {
		links = append(links, buildlink(node, i+1))
	}
	fmt.Printf("%+v", links)
	return links, nil
}

func linkParse(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkParse(c)...)
	}

	return ret
}

func buildlink(node *html.Node, id int) Link {
	var res Link

	res.ID = id
	for _, n := range node.Attr {
		if n.Key == "href" {
			res.Href = n.Val
			break
		}
	}

	res.Text = text(node)

	return res
}

func text(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	if node.Type != html.ElementNode {
		return ""
	}

	var res string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		res += text(c)
	}

	res = strings.Join(strings.Fields(res), " ")

	return res
}
