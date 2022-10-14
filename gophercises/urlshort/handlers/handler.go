package handlers

import (
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
)

type PathURL []map[string]string

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
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathsToUrls, err := parseYAMLtoMap(yml)
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		redirectPath, ok := pathsToUrls[path]
		if ok {
			log.Println("redirect path found", redirectPath)
			http.Redirect(w, r, redirectPath, http.StatusFound)
		} else {
			log.Println("redirect path NOT found")
			fallback.ServeHTTP(w, r)
		}

	}, nil
}

func parseYAMLtoMap(yml []byte) (map[string]string, error) {
	p := PathURL{}
	err := yaml.Unmarshal(yml, &p)
	if err != nil {
		return nil, err
	}
	m := getPathURLMap(p)
	return m, nil
}

func getPathURLMap(p PathURL) map[string]string {
	r := make(map[string]string)
	for _, m := range p {
		r[m["path"]] = m["url"]
	}
	return r
}
