package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"sync"
	"tickets/scrape"
	"time"
)

const Address = ":3000"
const TemplatePath = "template.html"

var tmpl = template.Must(template.ParseFiles(TemplatePath))
var from string
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
	flag.StringVar(&from, "from", "local", "download from web or from local")
	flag.Parse()
	// TODO clean up here
	var events []scrape.Event
	if from == "remote" {
		log.Println("scraping from remote ...")
		starHubEvents, err := scrape.GetSiteEvents("https://www.stubhub.com", "scrape/stubhub_event.json")
		if err != nil {
			log.Fatal(err)
		}
		viaGogoEvents, err := scrape.GetSiteEvents("https://www.viagogo.com", "scrape/viagogo_event.json")
		if err != nil {
			log.Fatal(err)
		}
		events = combineAndFilterEvents(starHubEvents, viaGogoEvents)
	} else {
		log.Println("loading from local storage ...")
		starHubEvents, err := scrape.LoadJsonToEvents("scrape/stubhub_event.json")
		if err != nil {
			log.Fatal(err)
		}

		viaGogoEvents, err := scrape.LoadJsonToEvents("scrape/viagogo_event.json")
		if err != nil {
			log.Fatal(err)
		}
		events = combineAndFilterEvents(starHubEvents, viaGogoEvents)
	}
	combinedEvents := NewCombinedEvents(events)

	go func(ce *CombinedEvents) {
		for {
			time.Sleep(60 * time.Second)
			stubHubEvents, err := scrape.GetSiteEvents("https://www.stubhub.com", "scrape/stubhub_event.json")
			if err != nil {
				log.Fatal(err)
			}

			viaGogoEvents, err := scrape.GetSiteEvents("https://www.viagogo.com", "scrape/viagogo_event.json")
			if err != nil {
				log.Fatal(err)
			}
			events = combineAndFilterEvents(stubHubEvents, viaGogoEvents)
			log.Printf("Total events got %d sleep for %d seconds\n", len(events), 60)

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
