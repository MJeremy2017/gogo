package main

import (
	"adventure/parser"
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

// TODO template caches
// TODO redirect when chapter key not exists
// TODO handle http internal server error
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
	chapterKey := getRegisteredChapterKey(r)
	chapter := story[chapterKey]

	err = tmpl.Execute(w, chapter)
	LogFatalIfErr(err)
}

func getRegisteredChapterKey(r *http.Request) string {
	defaultKey := "intro"
	key := r.URL.Path[1:]
	if len(key) == 0 || key == "home" {
		return defaultKey
	}
	if _, ok := story[key]; !ok {
		log.Printf("invalid chapter name %s", key)
		return defaultKey
	}
	return key
}

func LogFatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
