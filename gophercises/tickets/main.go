package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
	"tickets/scrape"
	"time"
)

const Address = ":3000"
const TemplatePath = "template.html"
const sleepSec = 60

var tmpl = template.Must(template.ParseFiles(TemplatePath))
var mu sync.Mutex

type CombinedEvents struct {
	Events []scrape.Event
}

func (c *CombinedEvents) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	log.Println("refreshing events", len(c.Events))
	err := tmpl.Execute(w, c.Events)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *CombinedEvents) UpdateEvents(e []scrape.Event) {
	mu.Lock()
	defer mu.Unlock()

	c.Events = scrape.SortEventTicketsByPrice(e)
}

func NewCombinedEvents(events []scrape.Event) CombinedEvents {
	res := scrape.SortEventTicketsByPrice(events)
	return CombinedEvents{Events: res}
}

func main() {
	var events []scrape.Event

	log.Println("loading from local storage ...")
	stubHubEvents, err := scrape.LoadJsonToEvents("scrape/stubhub_event.json")
	if err != nil {
		log.Fatal(err)
	}

	viaGogoEvents, err := scrape.LoadJsonToEvents("scrape/viagogo_event.json")
	if err != nil {
		log.Fatal(err)
	}
	events = combineAndFilterEvents(stubHubEvents, viaGogoEvents)
	combinedEvents := NewCombinedEvents(events)

	go func(ce *CombinedEvents) {
		for {
			time.Sleep(time.Duration(sleepSec) * time.Second)
			stubHubEvents, err := scrape.GetSiteEvents("https://www.stubhub.com", "scrape/stubhub_event.json")
			if err != nil {
				log.Fatal(err)
			}

			viaGogoEvents, err := scrape.GetSiteEvents("https://www.viagogo.com", "scrape/viagogo_event.json")
			if err != nil {
				log.Fatal(err)
			}
			events = combineAndFilterEvents(stubHubEvents, viaGogoEvents)
			log.Printf("Total events got %d sleep for %d seconds\n", len(events), sleepSec)

			ce.UpdateEvents(events)
		}
	}(&combinedEvents)

	mux := http.NewServeMux()
	mux.Handle("/", &combinedEvents)

	log.Println("listening on port", Address)
	log.Fatal(http.ListenAndServe(Address, mux))
}

func combineAndFilterEvents(eventList ...[]scrape.Event) []scrape.Event {
	var res []scrape.Event
	for _, events := range eventList {
		for _, event := range events {
			if len(event.Tickets) == 0 || event.Tickets[0].Price == 0 {
				continue
			}
			res = append(res, event)
		}
	}
	return res
}
