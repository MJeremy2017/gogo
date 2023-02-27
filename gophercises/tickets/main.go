package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"tickets/scrape"
)

const Address = ":3000"
const TemplatePath = "template.html"

var tmpl = template.Must(template.ParseFiles(TemplatePath))
var from string

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
	res := scrape.SortEventTicketsByPrice(viagogoEvents)
	return CombinedEvents{ViagogoEvents: res}
}

func main() {
	flag.StringVar(&from, "from", "local", "download from web or from local")
	flag.Parse()
	// TODO load and refresh the storage in the background
	var events []scrape.Event
	if from == "remote" {
		log.Println("scraping from remote ...")
		starHubEvents, err := scrape.GetSiteEvents("https://www.starhub.com", "scrape/starhub_event.json")
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
		starHubEvents, err := scrape.LoadJsonToEvents("scrape/starhub_event.json")
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

	mux := http.NewServeMux()
	mux.Handle("/", combinedEvents)

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
