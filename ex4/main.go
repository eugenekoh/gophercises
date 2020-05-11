package main

import (
	"fmt"
	"log"
	"os"

	"github.com/eugenekoh/gophercises/ex4/linkparser"
)

func main() {
	// constants
	const htmlFileName = "examples/ex2.html"

	f, err := os.Open(htmlFileName)
	if err != nil {
		log.Printf("unable to open html file\n %v\n", err)
	}

	links, err := linkparser.ParseLinks(f)
	if err != nil {
		log.Printf("error parsing html file\n %v\n", err)
	}

	fmt.Printf("%+v", links)
}
