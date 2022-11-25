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
const HtmlContent = `
<html>
    <head>
        <title>An adventure golang project</title>
        <style>
		  body {background-color: powderblue;}
		  h2 {color: black;}
		  p {color: blue;}
		</style>
    </head>
    <body>
        <h2>TITLE H2</h2>
        <p>This is the paragraph</p>
    </body>
</html>
`

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
	_, err := fmt.Fprint(w, HtmlContent)
	if err != nil {
		log.Fatal(err)
	}
}
