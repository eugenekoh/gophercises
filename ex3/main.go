package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	fileName := flag.String("-j", "gopher.json", "json file that contains story arcs")
	templatePath := flag.String("-t", "template.html", "html template")
	startArc := flag.String("-n", "intro", "default story arc name")

	flag.Parse()

	// initialize data structure
	storyArcs := parseJSON(*fileName)

	// load template
	tpl, err := template.ParseFiles(*templatePath)
	if err != nil {
		panic(err)
	}

	handler := cyoaHandler{
		storyArcs: storyArcs,
		tpl:       tpl,
		startArc:  *startArc,
	}

	http.Handle("/", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

// StoryArc is a data structure that contains info about current story arc
type StoryArc struct {
	Title   string    `json:"title"`
	Story   []string  `json:"story"`
	Options []Options `json:"options"`
}

// Options is a data structure that contains the options available for the current story arc
type Options struct {
	Text    string `json:"text"`
	NextArc string `json:"arc"`
}

func parseJSON(fileName string) map[string]StoryArc {
	jsonString, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	storyArcs := map[string]StoryArc{}

	json.Unmarshal(jsonString, &storyArcs)

	return storyArcs
}

type cyoaHandler struct {
	storyArcs map[string]StoryArc
	tpl       *template.Template
	startArc  string
}

func (h cyoaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	storyArcID := strings.TrimLeft(r.URL.Path, "/")
	storyArc, ok := h.storyArcs[storyArcID]
	fmt.Println(storyArcID)

	if !ok {
		storyArc = h.storyArcs[h.startArc]
	}

	fmt.Println(storyArc)
	h.tpl.Execute(w, storyArc)
}
