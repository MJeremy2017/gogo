package parser

import (
	"github.com/Workiva/go-datastructures/queue"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func BrowseLinks(url string) []string {
	var links []string
	parentUrl := url
	visited := make(map[string]bool)
	q := queue.New(1)
	err := q.Put(url)
	if err != nil {
		log.Fatalf("failed to put url %s %v", url, err)
	}
	for q.Len() > 0 {
		u, _ := q.Get(1)
		strUrl := u[0].(string)
		if !strings.Contains(strUrl, "http") && !strings.Contains(strUrl, "https") {
			strUrl = parentUrl + strUrl
		}
		visited[strUrl] = true

		urls, err := ParseUrlLinks(strUrl)
		if err != nil {
			log.Printf("error parsring url %s %v", strUrl, err)
			continue
		}
		links = append(links, strUrl)
		for _, u := range urls {
			if !visited[u] {
				_ = q.Put(u)
				visited[u] = true
			}
		}
	}

	return links
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
