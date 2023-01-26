package scrape_test

import (
	"fmt"
	"log"
	"testing"
	"tickets/scrape"
)

// TODO start from concert tickets and see how to find categories
func TestA(t *testing.T) {
	//depthFunc := colly.MaxDepth(2)
	//c := colly.NewCollector(depthFunc)
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
}
