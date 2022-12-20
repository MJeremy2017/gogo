package parser

import (
	"encoding/xml"
	"github.com/Workiva/go-datastructures/queue"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URL     []struct {
		Loc string `xml:"loc"`
	} `xml:"url"`
}

type UrlDepth struct {
	url   string
	depth int
}

func BrowseLinks(url string, maxDepth int) []string {
	var links []string
	parentUrl := url
	visited := make(map[string]struct{})

	q := queue.New(1)
	err := q.Put(UrlDepth{parentUrl, 1})
	if err != nil {
		log.Fatalf("failed to put url %s %v", url, err)
	}
	for q.Len() > 0 {
		u, _ := q.Get(1)
		urlDepth := u[0].(UrlDepth)
		baseUrl := urlDepth.url
		currDepth := urlDepth.depth
		if _, ok := visited[baseUrl]; ok {
			continue
		}
		if currDepth > maxDepth {
			continue
		}

		visited[baseUrl] = struct{}{}
		links = append(links, baseUrl)

		log.Println("processing", baseUrl)
		urls, err := ParseUrlLinks(baseUrl)
		if err != nil {
			log.Printf("error parsring url %s %v", baseUrl, err)
			continue
		}
		for _, u := range urls {
			fullPath := buildFullPath(baseUrl, u)
			if _, ok := visited[fullPath]; !ok {
				_ = q.Put(UrlDepth{fullPath, currDepth + 1})
			}
		}
	}

	return links
}

func EncodeLinksToXML(links []string, w io.Writer) error {
	var urlSet UrlSet
	for _, link := range links {
		urlSet.URL = append(urlSet.URL, struct {
			Loc string `xml:"loc"`
		}{Loc: link})
	}
	enc := xml.NewEncoder(w)
	enc.Indent("", "    ")
	if err := enc.Encode(urlSet); err != nil {
		return err
	}
	return nil
}

func buildFullPath(baseUrl string, path string) string {
	b := strings.TrimRight(baseUrl, "/")
	switch {
	case !isAbsPath(path):
		return b + formatRelativePath(path)
	case isAbsPath(path):
		return path
	}
	return path
}

func formatRelativePath(url string) string {
	return strings.TrimLeft(url, ".")
}

func isAbsPath(url string) bool {
	return strings.HasPrefix(url, "http")
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
