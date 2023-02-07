package main

import (
	"fmt"
	"log"
	"net/http"
	"tickets/scrape"
)

// TODO add a server
// TODO html page
// event name | time | venue | cheapest ticket (quantity price) | platform

type CombinedEvents struct {
	ViagogoEvents *[]scrape.Event
}

func (c CombinedEvents) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "Hello")
}

func NewCombinedEvents(viagogoEvents *[]scrape.Event) CombinedEvents {
	return CombinedEvents{ViagogoEvents: viagogoEvents}
}

func main() {
	events, err := scrapeViagogoTicket()
	if err != nil {
		log.Fatal(err)
	}
	combinedEvents := NewCombinedEvents(events)

	mux := http.NewServeMux()
	mux.Handle("/", combinedEvents)

	log.Fatal(http.ListenAndServe(":3000", mux))
}

func scrapeViagogoTicket() (*[]scrape.Event, error) {
	baseUrl := "https://www.viagogo.com"
	s := scrape.NewScraper(baseUrl)

	events, err := s.GetEvents("/sg/Concert-Tickets/Rock-and-Pop/Grace-Jones-Tickets")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for _, e := range events {
		err := s.GetTickets(&e)
		if err != nil {
			log.Println("Failed to get tickets info", e.TicketLink)
			continue
		}
		display(&e)
	}
	return &events, nil
}

func display(e *scrape.Event) {
	fmt.Printf("event name: %s, time: %s, venue %s, link: %s, ticket: %+v\n",
		e.EventName, e.Time, e.Venue, e.TicketLink, e.Tickets)
}
