package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	fileName := flag.String("-j", "gopher.json", "json file that contains story arcs")
	templatePath := flag.String("-t", "template.html", "html template")
	startArc := flag.String("-n", "intro", "default story arc name")

	flag.Parse()

	// initialize data structure
	story, err := parseJSON(*fileName)
	if err != nil {
		panic(err)
	}

	// load template
	tpl, err := template.ParseFiles(*templatePath)
	if err != nil {
		panic(err)
	}

	handler := cyoaHandler{
		story:    story,
		tpl:      tpl,
		startArc: *startArc,
	}

	http.Handle("/", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

// Story is a map collection of StoryArc
type Story map[string]StoryArc

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

func parseJSON(fileName string) (Story, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	story := Story{}

	d := json.NewDecoder(file)
	err = d.Decode(&story)

	if err != nil {
		return nil, err
	}

	return story, nil
}

type cyoaHandler struct {
	story    Story
	tpl      *template.Template
	startArc string
}

func (h cyoaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	storyArcID := strings.TrimLeft(r.URL.Path, "/")
	storyArc, ok := h.story[storyArcID]

	if !ok {
		storyArc = h.story[h.startArc]
	}

	.tpl.Execute(w, storyArc)
}
