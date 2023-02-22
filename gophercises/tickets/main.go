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

	var events []scrape.Event
	var err error
	if from == "remote" {
		log.Println("scraping from viagogo ...")
		events, err = scrapeViagogoTicket()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("loading from local storage ...")
		events, err = scrape.LoadJsonToEvents("scrape/viagogo_event.json")
		if err != nil {
			log.Fatal(err)
		}
	}
	combinedEvents := NewCombinedEvents(events)

	mux := http.NewServeMux()
	mux.Handle("/", combinedEvents)

	log.Println("listening on port", Address)
	log.Fatal(http.ListenAndServe(Address, mux))
}

func scrapeViagogoTicket() ([]scrape.Event, error) {
	p := "scrape/viagogo_event.json"
	s := scrape.NewScraper("https://www.viagogo.com")
	events := s.GetViagogoAllEvents()
	scrape.SaveEventsToJson(events, p)
	return events, nil
}
