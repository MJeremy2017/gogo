package scrape

import (
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

// FindLinks returns a map with key `category name` and value `relative path after host:port`
func (s *Scraper) FindLinks(path, query string) (map[string]string, error) {
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

// TODO in the response when the TicketsLeftInListingMessage is nil, means the ticket is sold
// TODO round floats
// Get all available tickets for an event
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
		q, ok := item["QuantityRange"].(string)
		if !ok {
			q = ""
		}
		p, ok := item["RawPrice"].(float64)
		if !ok {
			log.Printf("failed to convert raw price to %v float64\n", item["RawPrice"])
			p = 0.0
		}
		tickets = append(tickets, Ticket{
			QuantityRange: q,
			Price:         p,
		})
	}
	event.Tickets = tickets
	return nil
}
