package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Kotrice/gophercises/link"
)

func main() {
	filename := flag.String("html", "ex1.html", "html file name wanted to parse")
	flag.Parse()

	// open file
	file, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// parse html
	links, err := link.Parse(file)
	if err != nil {
		panic(err)
	}

	// print out result
	for _, link := range links {
		fmt.Printf("%+v\n", link)
	}
}
