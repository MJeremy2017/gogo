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
	fmt.Printf("\n got: %v", ms)
}
