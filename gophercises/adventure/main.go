package main

import (
	"adventure/parser"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const ADDRESS = ":8000"
const StoryFilePath = "story.json"
const HtmlTemplatePath = "story_template.html"
const StoryDefaultKey = "intro"

var tmpl = template.Must(template.ParseFiles(HtmlTemplatePath))

// TODO add test for story handler

func main() {
	mux := getRegisteredHandler()
	log.Println("listening on port", ADDRESS)
	log.Fatal(http.ListenAndServe(ADDRESS, mux))
}

func getRegisteredHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", storyHandler)
	return mux
}

func storyHandler(w http.ResponseWriter, r *http.Request) {
	story, err := parser.ParseStory(StoryFilePath)
	LogFatalIfErr(err)

	key := extractKey(r)
	chapter, err := getChapter(story, key)
	logAndRedirectWhenErr(w, r, err)

	err = tmpl.Execute(w, chapter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func extractKey(r *http.Request) string {
	key := r.URL.Path[1:]
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

func getChapter(story parser.Story, key string) (parser.Chapter, error) {
	if c, ok := story[key]; ok {
		return c, nil
	}
	return parser.Chapter{}, fmt.Errorf("invalid chapter key %s", key)
}

func LogFatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
