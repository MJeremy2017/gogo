package scrape_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"tickets/scrape"
)

// TODO save events to json
func TestA(t *testing.T) {
	s := scrape.NewScraper("https://www.viagogo.com")
	events := s.GetAllEvents()
	scrape.SaveEventsToJson(events)
}

func getAndSaveResponse(url string) {
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		log.Println(err)
	}

	bytes, _ := ioutil.ReadAll(response.Body)

	f, _ := os.Create("test.html")
	defer f.Close()
	f.WriteString(string(bytes))
	fmt.Println(string(bytes))
}

func postAndSaveResponse(url string) {
	response, err := http.Post(url, "application/json", nil)
	defer response.Body.Close()
	if err != nil {
		log.Println(err)
	}

	bytes, _ := ioutil.ReadAll(response.Body)

	f, _ := os.Create("test.json")
	defer f.Close()
	f.WriteString(string(bytes))
	fmt.Println(string(bytes))
}
