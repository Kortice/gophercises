package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Kotrice/gophercises/cyoa"
)

func main() {
	filename := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	port := flag.Int("port", 3030, "port the CYOA application on")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JSONStory(file)
	if err != nil {
		panic(err)
	}

	// WithTemplate test
	// t := template.Must(template.New("").Parse("hello world"))
	// handler := cyoa.NewHandler(story, cyoa.WithTemplate(t))

	handler := cyoa.NewHandler(story)

	fmt.Printf("Starting the sever on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}
