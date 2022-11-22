package main

import (
	"fmt"
	"log"
	"net/http"
)

func storyHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "hello")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", storyHandler)
	err := http.ListenAndServe("localhost:8000", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
