package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"tickets/scrape"
)

const Address = ":3000"
const TemplatePath = "template.html"

var tmpl = template.Must(template.ParseFiles(TemplatePath))

type CombinedEvents struct {
	ViagogoEvents []scrape.Event
}

func (c CombinedEvents) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := tmpl.Execute(w, c.ViagogoEvents)
	if err != nil {
		log.Fatal(err)
	}
}

func NewCombinedEvents(viagogoEvents []scrape.Event) CombinedEvents {
	return CombinedEvents{ViagogoEvents: viagogoEvents}
}

func main() {
	// TODO sort to put the cheapest price at the front
	events, err := scrapeViagogoTicket()
	if err != nil {
		log.Fatal(err)
	}
	combinedEvents := NewCombinedEvents(events)

	mux := http.NewServeMux()
	mux.Handle("/", combinedEvents)

	log.Println("listening on port", Address)
	log.Fatal(http.ListenAndServe(Address, mux))
}

func scrapeViagogoTicket() ([]scrape.Event, error) {
	baseUrl := "https://www.viagogo.com"
	s := scrape.NewScraper(baseUrl)

	events, err := s.GetEvents("/sg/Concert-Tickets/Rock-and-Pop/Grace-Jones-Tickets")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var result []scrape.Event
	for _, e := range events {
		err := s.GetTickets(&e)
		if err != nil {
			log.Println("Failed to get tickets info", e.TicketLink)
			continue
		}
		display(&e)
		result = append(result, e)
	}
	return result, nil
}

func display(e *scrape.Event) {
	fmt.Printf("event name: %s, time: %s, venue %s, link: %s, ticket: %+v\n",
		e.EventName, e.Time, e.Venue, e.TicketLink, e.Tickets)
}
