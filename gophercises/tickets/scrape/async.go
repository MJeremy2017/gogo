package scrape

import (
	"fmt"
	"log"
	"strings"
	"time"
)

var events []Event

func AsyncScrapeSiteEvents(host, fp string) {
	// TODO infinite loop
	doneChan := make(chan struct{})
	for {
		go func() {
			s := NewScraper(host)
			if strings.Contains(host, "stubhub") {
				events = s.GetStarHubAllEvents()
			} else if strings.Contains(host, "viagogo") {
				events = s.GetViagogoAllEvents()
			} else {
				log.Fatalln("unexpected host", host)
			}
			doneChan <- struct{}{}
		}()
		<-doneChan
		SaveEventsToJson(events, fp)
		fmt.Println("sleep for 5 seconds ...")
		time.Sleep(time.Second * 5)
	}
}
