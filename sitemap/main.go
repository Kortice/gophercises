package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Kotrice/gophercises/link"
)

const XMLNS = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url you want to build sitemap for")
	maxDepth := flag.Int("depth", 3, "the maximum number of depth traverse for")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)

	toXml := urlset{
		Urls:  make([]loc, len(pages)),
		Xmlns: XMLNS,
	}

	fmt.Print(xml.Header)
	encoder := xml.NewEncoder(os.Stdout)
	defer encoder.Close()
	encoder.Indent("", "  ")

	for i, page := range pages {
		toXml.Urls[i] = loc{page}
	}

	if err := encoder.Encode(toXml); err != nil {
		panic(err)
	}

}

func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{})

	var q []string
	q = append(q, urlStr)
	seen[urlStr] = struct{}{}

	step := 0

	for len(q) > 0 {
		length := len(q)
		for range length {
			curUrl := q[0]
			q = q[1:]

			for _, l := range get(curUrl) {
				if _, ok := seen[l]; !ok {
					q = append(q, l)
					seen[l] = struct{}{}
				}
			}
		}
		step++
		if step == maxDepth {
			break
		}
	}

	ret := make([]string, 0, len(seen))
	for k := range seen {
		ret = append(ret, k)
	}

	return ret
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()

	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, err := link.Parse(r)
	if err != nil {
		panic(err)
	}

	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}

	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string

	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}

	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}
