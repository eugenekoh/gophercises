package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/eugenekoh/gophercises/ex2/urlshort"
)

func main() {
	jsonFileName := flag.String("json", "data.json", "a json file to specify paths")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler, err := urlshort.MapHandler(pathsToUrls, mux)
	if err != nil {
		panic(err)
	}

	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	jsonFile, err := ioutil.ReadFile(*jsonFileName)
	if err != nil {
		panic(err)
	}

	// Build the JSONHandler using the YAMLHandler as the
	// fallbac
	jsonHandler, err := urlshort.JSONHandler(jsonFile, yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
