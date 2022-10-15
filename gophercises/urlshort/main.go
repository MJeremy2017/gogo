package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"urlshort/handlers"
)

var yamlFile string
var jsonFile string

func main() {
	parseFlags()
	mux := defaultMux()

	jsonHandler, err := handlers.JSONHandler(jsonFile, mux)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	err = http.ListenAndServe(":8080", jsonHandler)
	if err != nil {
		log.Fatalln("unable to start the server", err)
	}
}

func parseFlags() {
	flag.StringVar(&yamlFile, "yml", "", "yaml file name to stored urls")
	flag.StringVar(&jsonFile, "json", "", "json file name to stored urls")
	flag.Parse()
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "Hello, world!")
}
