package cyoa

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func init() {
	tmpl = template.Must(template.ParseFiles("../../template.html"))
}

var tmpl *template.Template

type Story map[string]Arc

type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tmpl, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)
	if arc, ok := h.s[path]; ok {
		err := h.t.Execute(w, arc)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Story arc not found...", http.StatusNotFound)
}

func ParseJson(data []byte) (Story, error) {
	var storyData Story
	if err := json.Unmarshal(data, &storyData); err != nil {
		return nil, err
	}
	return storyData, nil
}

func (s Story) PrintArcTitles() {
	for k := range s {
		fmt.Printf("%s: %s\n", k, s[k].Title)
	}
}
