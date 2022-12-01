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

var story parser.Story
var err error
var tmpl = template.Must(template.ParseFiles(HtmlTemplatePath))

func main() {
	story, err = parser.ParseStory(StoryFilePath)
	LogFatalIfErr(err)

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
	chapter, err := getChapter(r)
	logAndRedirectWhenErr(w, r, err)
	err = tmpl.Execute(w, chapter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func logAndRedirectWhenErr(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		errMsg := fmt.Sprintf("unexpected error when directing chapters %+v", err)
		log.Println(errMsg)
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func getChapter(r *http.Request) (parser.Chapter, error) {
	defaultKey := "intro"
	key := r.URL.Path[1:]
	if len(key) == 0 || key == "home" {
		return story[defaultKey], nil
	}
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
