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
	// /sg/Concert-Tickets
	categories, err := s.FindLinks("", scrape.CategoryQuery)
	if err != nil {
		log.Println(err)
	}

	for _, catLink := range categories {
		// /sg/Concert-Tickets/Clubs-and-Dance
		eventTypes, err := s.FindLinks(catLink, scrape.EventTypeQuery)
		if err != nil {
			log.Println(err)
		}

		for _, etLink := range eventTypes {
			// /sg/Concert-Tickets/Rock-and-Pop/Bastille-Tickets
			eventLinks, err := s.FindLinks(etLink, scrape.EventQuery)
			if err != nil {
				log.Println(err)
			}

			for _, eL := range eventLinks {
				// event object
				events, err := s.GetEvents(eL)
				if err != nil {
					log.Println(err)
				}
				for i, e := range events {
					err := s.GetTickets(&e)
					if err != nil {
						log.Println("Failed to get tickets info -------")
						continue
					}
					if len(e.Tickets) == 0 {
						continue
					}
					fmt.Printf("Event %d name: %s, time: %s, venue %s, link: %s, ticket: %+v\n",
						i+1, e.EventName, e.Time, e.Venue, e.TicketLink, e.Tickets)

					os.Exit(0)
				}
			}
		}
	}
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
