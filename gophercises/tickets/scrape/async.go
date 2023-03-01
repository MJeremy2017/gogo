package scrape

import (
	"log"
	"strings"
)

var events []Event

func AsyncScrapeSiteEvents(host, fp string) {
	// TODO infinite loop
	chanEvents := make(chan []Event, 1)
	for {
		doneChan := make(chan bool, 1)
		go func() {
			s := NewScraper(host)
			if strings.Contains(host, "stubhub") {
				events = s.GetStarHubAllEvents()
			} else if strings.Contains(host, "viagogo") {
				events = s.GetViagogoAllEvents()
			} else {
				log.Fatalln("unexpected host", host)
			}
			chanEvents <- events
			doneChan <- true
		}()
		events = <-chanEvents
		SaveEventsToJson(events, fp)

	}
}
