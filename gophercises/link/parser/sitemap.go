package parser

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func BrowseLinks(url string) {
	html, _ := GetHtmlPage(url)
	fmt.Println(html)
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
