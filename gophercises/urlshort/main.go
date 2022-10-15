package main

import (
	"flag"
	"fmt"
	"github.com/go-errors/errors"
	"log"
	"net/http"
	"urlshort/handlers"
)

var yamlFile string
var jsonFile string

func main() {
	parseFlags()
	mux := defaultMux()

	handler, err := handlers.BoltDbHandler(mux)
	if err != nil {
		log.Fatalln(err.(*errors.Error).ErrorStack())
	}
	fmt.Println("Starting the server on :8080")
	err = http.ListenAndServe(":8080", handler)
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
