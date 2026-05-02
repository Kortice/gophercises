package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Link struct {
	Href string
	Text string
}

func main() {
	filename := flag.String("html", "ex1.html", "html file name wanted to parse")
	flag.Parse()

	// open file
	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal("open file failed", err)
	}
	defer file.Close()

	// parse html
	doc, err := html.Parse(file)
	if err != nil {
		log.Fatal("parse html failed", err)
	}

	links := HTMLLinkParser(doc)

	for _, link := range links {
		fmt.Printf("+%v\n", link)
	}
}

func HTMLLinkParser(doc *html.Node) (links []Link) {
	for n := range doc.Descendants() {
		// when <a> </a>
		if n.Type == html.ElementNode && n.DataAtom == atom.A {

			link := Link{}

			// set link.Href
			for _, a := range n.Attr {
				if a.Key == "href" {
					link.Href = a.Val
					break
				}
			}

			// set link.Text
			var text []string
			for sub := range n.Descendants() {
				if sub.Type == html.TextNode {
					text = append(text, strings.TrimSpace(sub.Data))
				}
			}
			link.Text = strings.Join(text, " ")

			// append to links
			links = append(links, link)
		}
	}

	return
}
