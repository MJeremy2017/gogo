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
	Platform   string
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

func (s *Scraper) GetViaGogoEvents(path string) ([]Event, error) {
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
			Platform:   "viagogo",
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
		e.Platform = "Viagogo"
		events = append(events, *e)
	}

	return events, nil
}

func (s *Scraper) GetStarHubEvents(path string) ([]Event, error) {
	sc := struct {
		CategorySummary map[string]interface{} `json:"categorySummary"`
		EventGrids      struct {
			Num2 struct {
				Items []map[string]interface{} `json:"items"`
			} `json:"2"`
		} `json:"eventGrids"`
	}{}
	var res []Event

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
	eventName := getStringFromMap(sc.CategorySummary, "categoryName")
	for _, it := range sc.EventGrids.Num2.Items {
		date := getStringFromMap(it, "formattedDate")
		hour := getStringFromMap(it, "formattedTime")
		v1 := getStringFromMap(it, "venueName")
		v2 := getStringFromMap(it, "formattedVenueLocation")
		minPrice := getStringFromMap(it, "formattedMinPrice")
		res = append(res, Event{
			EventName:  eventName,
			Time:       formatDateHour(date, hour),
			Venue:      v1 + ";" + v2,
			TicketLink: getStringFromMap(it, "url"),
			Tickets: []Ticket{
				{
					Price:  formatDollarSignPrice(minPrice),
					BuyUrl: getStringFromMap(it, "url"),
				},
			},
			Platform: "StarHub",
		})
	}
	return res, nil
}

func (s *Scraper) GetStarHubAllEvents() []Event {
	concertPath := "/concert-tickets/category/1/"
	links, err := s.FindStarHubEventLinks(concertPath)
	if err != nil {
		log.Println(err)
		return nil
	}
	var res []Event
	for _, lk := range links {
		es, err := s.GetStarHubEvents(lk)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, e := range es {
			if e.Tickets[0].Price == 0 {
				// no ticket left
				continue
			}
			u := e.Tickets[0].BuyUrl
			fullUrl := s.joinPath(s.baseUrl, u)
			e.Tickets[0].BuyUrl = fullUrl
			display(&e)
			res = append(res, e)
		}
	}
	return res
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
		q := getStringFromMap(item, "QuantityRange")
		p := getFloatFromMap(item)
		u := getStringFromMap(item, "BuyUrl")
		tickets = append(tickets, Ticket{
			QuantityRange: q,
			Price:         RoundRawPrice(p),
			BuyUrl:        s.joinPath(s.baseUrl, u),
		})
	}
	event.Tickets = tickets
	return nil
}

func (s *Scraper) GetViagogoAllEvents() []Event {
	var result []Event
	var cnt int
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
				events, err := s.GetViaGogoEvents(eL)
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
					result = append(result, e)
					cnt += 1
					if cnt%10 == 0 {
						fmt.Println("processed", cnt)
						display(&e)
					}
				}
			}
		}
	}
	return result
}

func GetSiteEvents(host, fp string) ([]Event, error) {
	var events []Event
	s := NewScraper(host)
	if host == "https://www.stubhub.com" {
		events = s.GetStarHubAllEvents()
	} else {
		events = s.GetViagogoAllEvents()
	}
	SaveEventsToJson(events, fp)
	return events, nil
}

func display(e *Event) {
	fmt.Printf("event name: %s, time: %s, venue %s, link: %s, ticket: %+v\n",
		e.EventName, e.Time, e.Venue, e.TicketLink, e.Tickets)
}
