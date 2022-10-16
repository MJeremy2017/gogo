package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"urlshort/data"
)

type PathURL []map[string]string

const defaultYaml = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

// MapHandler will return a http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		redirectPath, ok := pathsToUrls[path]
		if ok {
			log.Println("redirect path found", redirectPath)
			http.Redirect(writer, request, redirectPath, http.StatusFound)
		} else {
			log.Println("redirect path NOT found")
			fallback.ServeHTTP(writer, request)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// a http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(file string, fallback http.Handler) (http.HandlerFunc, error) {
	var yml []byte
	yml, err := readFromFile(file)
	if err != nil {
		log.Printf("unable to load yaml from file %v", err)
		yml = []byte(defaultYaml)
	}
	pathsToUrls, err := parseYAMLtoMap(yml)
	if err != nil {
		panic(err)
	}
	return MapHandler(pathsToUrls, fallback), nil
}

// JSONHandler will handler and path to urls from a given json file
func JSONHandler(file string, fallback http.Handler) (http.HandlerFunc, error) {
	var js []byte
	js, err := readFromFile(file)
	if err != nil {
		panic(err)
	}
	pathsToUrls, err := parseJSONtoMap(js)
	if err != nil {
		panic(err)
	}
	return MapHandler(pathsToUrls, fallback), nil
}

// BoltDbHandler saves and accesses data from bolt DB
func BoltDbHandler(fallback http.Handler) (http.HandlerFunc, error) {
	db := data.NewBoltDB("path_to_urls")
	err := writeDataToBoltdb(db)
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		redirectPath := db.Get([]byte(p))
		if redirectPath != "" {
			log.Println("redirect path found", redirectPath)
			http.Redirect(w, r, redirectPath, http.StatusFound)
		} else {
			log.Println("redirect path NOT found")
			fallback.ServeHTTP(w, r)
		}
	}, nil
}

func writeDataToBoltdb(db *data.BoltDB) error {
	pathToUrls := map[string]string{
		"/urlshort": "https://github.com/gophercises/urlshort",
		"/bolt":     "https://github.com/boltdb/bolt#getting-started",
	}
	for k, v := range pathToUrls {
		err := db.Set([]byte(k), []byte(v))
		if err != nil {
			return errors.Wrap(err, 0)
		}
		fmt.Println("loaded data", k, v)
	}
	return nil
}

func parseJSONtoMap(js []byte) (map[string]string, error) {
	p := PathURL{}
	err := json.Unmarshal(js, &p)
	if err != nil {
		return nil, err
	}
	return getPathURLMap(p), nil
}

func readFromFile(file string) ([]byte, error) {
	f, err := os.ReadFile(file)
	return f, err
}

func parseYAMLtoMap(yml []byte) (map[string]string, error) {
	p := PathURL{}
	err := yaml.Unmarshal(yml, &p)
	if err != nil {
		return nil, err
	}
	return getPathURLMap(p), nil
}

func getPathURLMap(p PathURL) map[string]string {
	r := make(map[string]string)
	for _, m := range p {
		r[m["path"]] = m["url"]
	}
	return r
}
