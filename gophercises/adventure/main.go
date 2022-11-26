package main

import (
	"adventure/parser"
	"html/template"
	"log"
	"net/http"
)

// TODO: Add links to '/'
// TODO: for each story, render a template (one method for all)

const ADDRESS = ":8000"
const StoryFilePath = "story.json"
const HtmlTemplatePath = "story_template.html"

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

	tmpl, err := template.ParseFiles(HtmlTemplatePath)
	LogFatalIfErr(err)

	err = tmpl.Execute(w, story)
	LogFatalIfErr(err)
}

func LogFatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
