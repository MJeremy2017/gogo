package main

import (
	"adventure/parser"
	"html/template"
	"log"
	"net/http"
)

// TODO: for each story, render a template (one method for all)
// {story_name: {title, story, options: {text, arc_links}}}
// host:port/story_name -> render template

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

	chapter := story["intro"]
	tmpl, err := template.ParseFiles(HtmlTemplatePath)
	LogFatalIfErr(err)

	err = tmpl.Execute(w, chapter)
	LogFatalIfErr(err)
}

func LogFatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
