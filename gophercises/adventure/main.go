package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

const ADDRESS = ":8000"
const StoryFilePath = "story.json"
const HtmlTemplatePath = "story_template.html"
const StoryDefaultKey = "intro"

var tmpl = template.Must(template.ParseFiles(HtmlTemplatePath))

type Story map[string]Chapter

func (s Story) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := extractKey(r)
	chapter, err := s.getChapter(key)
	logAndRedirectWhenErr(w, r, err)

	err = tmpl.Execute(w, chapter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

func ParseStory(fp string) (Story, error) {
	b, err := os.ReadFile(fp)
	if err != nil {
		return Story{}, err
	}

	var story Story
	err = json.Unmarshal(b, &story)
	if err != nil {
		return Story{}, err
	}
	log.Println("story loaded successfully")
	return story, nil
}

func main() {
	story, err := ParseStory(StoryFilePath)
	LogFatalIfErr(err)

	mux := http.NewServeMux()
	mux.Handle("/", story)

	log.Println("listening on port", ADDRESS)
	log.Fatal(http.ListenAndServe(ADDRESS, mux))
}

func extractKey(r *http.Request) string {
	var key string
	key = strings.TrimLeft(r.URL.Path, "/")
	if len(key) == 0 || key == "home" {
		return StoryDefaultKey
	}
	return key
}

func logAndRedirectWhenErr(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		errMsg := fmt.Sprintf("unexpected error when directing chapters %+v", err)
		log.Println(errMsg)
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func (s Story) getChapter(key string) (Chapter, error) {
	if c, ok := s[key]; ok {
		return c, nil
	}
	return Chapter{}, fmt.Errorf("invalid chapter key %s", key)
}

func LogFatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
