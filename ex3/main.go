package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	//flags
	fileName := flag.String("-j", "gopher.json", "json file that contains story arcs")
	templatePath := flag.String("-t", "template.html", "html template")
	startArc := flag.String("-n", "intro", "default story arc name")
	port := flag.Int("-p", 3000, "port to serve server")

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

	templateOption := WithTemplate(tpl)

	handler := NewCyoaHandler(story, *startArc, templateOption)
	http.Handle("/", handler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))

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

// HandlerOption defines an option type for NewCyoaHandler
type HandlerOption func(h *cyoaHandler)

// WithTemplate is a function that adds a defined template to the handler function
func WithTemplate(tpl *template.Template) HandlerOption {
	return func(h *cyoaHandler) {
		h.tpl = tpl
	}
}

// NewCyoaHandler is a constructor for cyoaHandler
func NewCyoaHandler(s Story, startArc string, opts ...HandlerOption) http.Handler {

	// defaults
	tpl := template.Must(template.New("").Parse("Hello!"))
	pathFnc := func(r *http.Request) string {
		s := r.URL.Path
		return strings.TrimLeft(s, "/")
	}

	h := cyoaHandler{s, tpl, startArc, pathFnc}

	for _, opt := range opts {
		opt(&h)
	}

	return h
}

type cyoaHandler struct {
	story    Story
	tpl      *template.Template
	startArc string
	pathFnc  func(r *http.Request) string
}

func (h cyoaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	storyArcID := h.pathFnc(r)
	storyArc, ok := h.story[storyArcID]

	// no path found, default to startArc
	if !ok {
		storyArc = h.story[h.startArc]
	}

	log.Printf("%s", storyArcID)
	log.Printf("%+v", storyArc)
	err := h.tpl.Execute(w, storyArc)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
	}
}
