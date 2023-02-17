package scrape

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strings"
)

const CategoryQuery = ".prinav a[href]"
const EventTypeQuery = ".cloud a[href]"
const EventQuery = "div.uuxxl.pgw ul.cloud.mbxl a[href]"

type ticketItems struct {
	Items []map[string]interface{} `json:"Items"`
}

type Ticket struct {
	QuantityRange string
	Price         float64
	BuyUrl        string
}

type Event struct {
	EventName  string
	Time       string
	Venue      string
	TicketLink string
	Tickets    []Ticket
}

type Scraper struct {
	baseUrl string
}

func NewScraper(baseUrl string) *Scraper {
	return &Scraper{baseUrl: baseUrl}
}

// FindViaGogoLinks returns a map with key `category name` and value `relative path after host:port`
func (s *Scraper) FindViaGogoLinks(path, query string) (map[string]string, error) {
	res := make(map[string]string)

	c := colly.NewCollector()
	c.OnHTML(query, func(e *colly.HTMLElement) {
		relativePath := e.Attr("href")
		text := strings.TrimSpace(e.Text)
		res[text] = relativePath
	})

	url := s.joinPath(s.baseUrl, path)
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Scraper) FindStarHubEventLinks(path string) (map[string]string, error) {
	sc := struct {
		CategoryGridLinks []struct {
			ID   int    `json:"id"`
			Text string `json:"text"`
			URL  string `json:"url"`
		} `json:"categoryGridLinks"`
	}{}
	res := make(map[string]string)

	q := "script#index-data"
	c := colly.NewCollector()
	c.OnHTML(q, func(e *colly.HTMLElement) {
		err := json.Unmarshal([]byte(e.Text), &sc)
		if err != nil {
			log.Println(err)
			return
		}
	})

	url := s.joinPath(s.baseUrl, path)
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}
	for _, c := range sc.CategoryGridLinks {
		res[c.Text] = c.URL
	}
	return res, nil
}

func (s *Scraper) joinPath(baseUrl, path string) string {
	p := strings.TrimLeft(path, "/")
	return strings.TrimRight(baseUrl, "/") + "/" + p
}

func (s *Scraper) GetEvents(path string) ([]Event, error) {
	ms := make([]*Event, 0)
	var eventName string
	c := colly.NewCollector()
	c.OnHTML("#catNameInHeader", func(e *colly.HTMLElement) {
		if eventName == "" {
			eventName = strings.TrimSpace(e.Text)
		}
	})

	c.OnHTML(".js-event-row-container.el-row-anchor", func(e *colly.HTMLElement) {
		p := e.Attr("href")
		t := e.ChildAttr("time", "datetime")
		v := e.ChildText("span[class=t-b]")
		event := Event{
			Time:       t,
			TicketLink: p,
			Venue:      v,
		}
		ms = append(ms, &event)
	})

	url := s.joinPath(s.baseUrl, path)
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	var events []Event
	for _, e := range ms {
		e.EventName = eventName
		events = append(events, *e)
	}

	return events, nil
}

// GetTickets returns all available tickets for an event
func (s *Scraper) GetTickets(event *Event) error {
	var tickets []Ticket
	lk := event.TicketLink
	url := s.joinPath(s.baseUrl, lk)
	items, err := postAndGetJsonResponse(url)
	if err != nil {
		log.Println(err)
		return err
	}
	for _, item := range items.Items {
		if item["TicketsLeftInListingMessage"] == nil {
			continue
		}
		q := getQuantityRangeFromItem(item)
		p := getRawPriceFromItem(item)
		u := getBuyUrlFromItem(item)
		tickets = append(tickets, Ticket{
			QuantityRange: q,
			Price:         RoundRawPrice(p),
			BuyUrl:        s.joinPath(s.baseUrl, u),
		})
	}
	event.Tickets = tickets
	return nil
}

func (s *Scraper) GetAllEvents() []Event {
	var result []Event
	categories := map[string]string{
		"Concert Tickets": "/sg/Concert-Tickets",
	}

	for _, catLink := range categories {
		// /sg/Concert-Tickets/Clubs-and-Dance
		eventTypes, err := s.FindViaGogoLinks(catLink, EventTypeQuery)
		if err != nil {
			log.Println(err)
		}

		for _, etLink := range eventTypes {
			// /sg/Concert-Tickets/Rock-and-Pop/Bastille-Tickets
			eventLinks, err := s.FindViaGogoLinks(etLink, EventQuery)
			if err != nil {
				log.Println(err)
			}

			for _, eL := range eventLinks {
				// event object
				events, err := s.GetEvents(eL)
				if err != nil {
					log.Println(err)
				}
				for _, e := range events {
					err := s.GetTickets(&e)
					if err != nil {
						log.Println("Failed to get tickets info -------")
						continue
					}
					if len(e.Tickets) == 0 {
						continue
					}
					display(&e)
					result = append(result, e)
				}
			}
		}
	}
	return result
}

func display(e *Event) {
	fmt.Printf("event name: %s, time: %s, venue %s, link: %s, ticket: %+v\n",
		e.EventName, e.Time, e.Venue, e.TicketLink, e.Tickets)
}
