package parser

import (
	"fmt"
	"github.com/Workiva/go-datastructures/queue"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type UrlDepth struct {
	url   string
	depth int
}

func BrowseLinks(url string, maxDepth int) []string {
	var links []string
	parentUrl := url
	visited := make(map[string]bool)

	q := queue.New(1)
	err := q.Put(UrlDepth{parentUrl, 1})
	if err != nil {
		log.Fatalf("failed to put url s% %v", url, err)
	}
	for q.Len() > 0 {
		u, _ := q.Get(1)
		urlDepth := u[0].(UrlDepth)
		fullPath := urlDepth.url
		currDepth := urlDepth.depth
		if visited[fullPath] || currDepth > maxDepth {
			continue
		}
		visited[fullPath] = true

		links = append(links, fullPath)

		fmt.Println("processing", fullPath)
		urls, err := ParseUrlLinks(fullPath)
		if err != nil {
			log.Printf("error parsring url %s %v", fullPath, err)
			continue
		}
		for _, u := range urls {
			urlPath := formatUrl(u)
			fullPath := buildFullPath(parentUrl, urlPath)
			if !visited[fullPath] {
				_ = q.Put(UrlDepth{fullPath, currDepth + 1})
			}
		}
	}

	return links
}

func buildFullPath(baseUrl string, path string) string {
	if isNotValidUrl(path) {
		return baseUrl + "/" + path
	}
	return path
}

func formatUrl(url string) string {
	if isNotValidUrl(url) {
		url = strings.TrimLeft(strings.TrimLeft(url, "."), "/")
	}
	return url
}

func isNotValidUrl(url string) bool {
	return !strings.Contains(url, "http") && !strings.Contains(url, "https")
}

func ParseUrlLinks(url string) ([]string, error) {
	data, err := GetHtmlPage(url)
	if err != nil {
		return nil, err
	}
	p := NewParser(strings.NewReader(data))
	links, err := p.ParseLinks()

	if err != nil {
		return nil, err
	}

	var urls []string
	for _, lk := range links {
		urls = append(urls, lk.Href)
	}
	return urls, nil
}

func GetHtmlPage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
