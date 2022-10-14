package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"urlshort/handlers"
)

var yamlFile string

func main() {
	parseFlags()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := handlers.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := handlers.YAMLHandler(yamlFile, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	err = http.ListenAndServe(":8080", yamlHandler)
	if err != nil {
		log.Fatalln("unable to start the server", err)
	}
}

func parseFlags() {
	flag.StringVar(&yamlFile, "yml", "", "yaml file name to stored urls")
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
