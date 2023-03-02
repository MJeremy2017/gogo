package scrape_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"tickets/scrape"
	"time"
)

func TestViaGogo(t *testing.T) {
	events, err := scrape.LoadJsonToEvents("viagogo_event.json")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(events[0])
}

func TestStarHub(t *testing.T) {
	events, err := scrape.LoadJsonToEvents("stubhub_event.json")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(events[0])
}

func TestAsyncFetch(t *testing.T) {
	//scrape.AsyncScrapeSiteEvents("https://www.stubhub.com", "stubhub_event.json")
	a := 1
	go func() {
		b := 1
		for {
			fmt.Println("inside", b)
			b += 1
			time.Sleep(time.Second)
		}
	}()

	fmt.Println("outside", a)
	time.Sleep(10 * time.Second)

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
