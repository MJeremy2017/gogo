package scrape

import (
	"github.com/gocolly/colly"
	"strings"
)

const CategoryQuery = ".prinav a[href]"
const EventTypeQuery = ".cloud a[href]"
const EventQuery = "div.uuxxl.pgw ul.cloud.mbxl a[href]"

type Ticket struct {
	Quantity int32
	Price    float64
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

// TODO use post request to get the tickets info
// TODO in the response when the TicketsLeftInListingMessage is null, means the ticket is sold
func (s *Scraper) GetTickets(events []Event) error {
	return nil
}
