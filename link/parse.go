package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	nodes := linkNodes(doc)

	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}

	return links, nil
}

func buildLink(n *html.Node) Link {
	ret := Link{}

	// set Href
	for _, a := range n.Attr {
		if a.Key == "href" {
			ret.Href = a.Val
		}
	}

	// set Text
	ret.Text = text(n)

	return ret
}

func linkNodes(doc *html.Node) []*html.Node {
	var ret []*html.Node

	for n := range doc.Descendants() {
		// if <a>
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			ret = append(ret, n)
		}
	}

	return ret
}

func text(n *html.Node) string {
	var ret string

	for c := range n.Descendants() {
		if c.Type == html.TextNode {
			ret += c.Data
		}
	}

	return strings.Join(strings.Fields(ret), " ")
}
