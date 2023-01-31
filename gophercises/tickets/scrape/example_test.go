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

func TestA(t *testing.T) {
	s := scrape.NewScraper("https://www.viagogo.com")
	ms, err := s.FindLinks("sg/Concert-Tickets/Rock-and-Pop", scrape.EventQuery)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("\ngot: %v\n", ms)

	events, err := s.GetEvents("/sg/Concert-Tickets/Rock-and-Pop/Grace-Jones-Tickets")
	if err != nil {
		log.Println(err)
	}
	for i, e := range events {
		fmt.Printf("Event %d name: %s, time: %s, venue %s, link: %s \n", i+1, e.EventName, e.Time, e.Venue, e.TicketLink)
	}

	ticketUrl := "https://www.viagogo.com/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets/E-151336327"
	postAndSaveResponse(ticketUrl)
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
