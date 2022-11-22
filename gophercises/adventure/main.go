package main

import (
	"fmt"
	"log"
	"net/http"
)

// TODO: display an html format with links below
// TODO: use a html template to display
// TODO: load json in structs and put in template

const ADDRESS = ":8000"

func main() {
	mux := getRegisteredHandler()
	log.Println("listening on port", ADDRESS)
	log.Fatalln(http.ListenAndServe(ADDRESS, mux))
}

func getRegisteredHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", storyHandler)
	return mux
}

func storyHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "hello")
	if err != nil {
		log.Fatal(err)
	}
}
