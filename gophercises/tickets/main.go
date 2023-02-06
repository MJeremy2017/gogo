package main

import (
	"fmt"
	"log"
	"tickets/scrape"
)

// TODO add a server
// TODO html page
// event name | time | venue | cheapest ticket (quantity price) | platform

func main() {
	baseUrl := "https://www.viagogo.com"
	s := scrape.NewScraper(baseUrl)

	events, err := s.GetEvents("/sg/Concert-Tickets/Rock-and-Pop/Grace-Jones-Tickets")
	if err != nil {
		log.Println(err)
	}
	for _, e := range events {
		err := s.GetTickets(&e)
		if err != nil {
			log.Println("Failed to get tickets info", e.TicketLink)
			continue
		}
		display(&e)
	}
}

func display(e *scrape.Event) {
	fmt.Printf("event name: %s, time: %s, venue %s, link: %s, ticket: %+v\n",
		e.EventName, e.Time, e.Venue, e.TicketLink, e.Tickets)
}
