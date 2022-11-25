package main

import (
	"adventure/parser"
	"html/template"
	"log"
	"net/http"
)

// TODO: display an html format with links below
// TODO: use a html template to display
// TODO: load json in structs and put in template

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
	if err != nil {
		log.Fatal(err)
	}
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
