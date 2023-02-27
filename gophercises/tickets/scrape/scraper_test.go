package scrape

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScraper_FindLinks(t *testing.T) {
	s := setUp()
	defer s.Close()
	scraper := NewScraper(s.URL)

	t.Run("Can find category links", func(t *testing.T) {
		want := map[string]string{
			"Concert Tickets":  "/sg/Concert-Tickets",
			"Sports Tickets":   "/sg/Sports-Tickets",
			"Theatre Tickets":  "/sg/Theatre-Tickets",
			"Festival Tickets": "/sg/Festival-Tickets",
		}
		got, err := scraper.FindViaGogoLinks("", CategoryQuery)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Can find event type links", func(t *testing.T) {
		want := map[string]string{
			"Club and dance": "/sg/Concert-Tickets/Clubs-and-Dance",
			"Flamenco":       "/sg/Theater-Tickets/Flamenco",
		}
		got, err := scraper.FindViaGogoLinks("/sg/Concert-Tickets", EventTypeQuery)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Can find events links", func(t *testing.T) {
		want := map[string]string{
			"5 Seconds of Summer": "/sg/Concert-Tickets/Rock-and-Pop/5-Seconds-of-Summer-Tickets",
			"Bastille":            "/sg/Concert-Tickets/Rock-and-Pop/Bastille-Tickets",
			"Brigitte":            "/sg/Concert-Tickets/Rock-and-Pop/Brigitte-Tickets",
		}
		got, err := scraper.FindViaGogoLinks("/sg/Concert-Tickets/Clubs-and-Dance", EventQuery)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("return empty when no events", func(t *testing.T) {
		want := make(map[string]string)
		got, err := scraper.FindViaGogoLinks("/", EventQuery)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("return all events ticket info", func(t *testing.T) {
		eventLink := "/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets"
		want := []Event{
			{
				EventName:  "Super Junior Tickets",
				Time:       "2023-02-09T20:00:00",
				Venue:      "Sao Paulo, Brazil",
				TicketLink: "/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets/E-151336327",
				Platform:   "Viagogo",
			},
			{
				EventName:  "Super Junior Tickets",
				Time:       "2023-02-18T19:30:00",
				Venue:      "Selangor, Malaysia",
				TicketLink: "/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets/E-151396140",
				Platform:   "Viagogo",
			},
		}
		got, err := scraper.GetViaGogoEvents(eventLink)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("return all tickets with prices", func(t *testing.T) {
		event := Event{
			EventName:  "Super Junior Tickets",
			Time:       "2023-02-09T20:00:00",
			Venue:      "Sao Paulo, Brazil",
			TicketLink: "/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets/E-151336327",
		}

		err := scraper.GetTickets(&event)
		assert.NoError(t, err)

		var got []Ticket
		for _, ticket := range event.Tickets {
			got = append(got, ticket)
		}

		want := []Ticket{
			{
				QuantityRange: "1 - 4",
				Price:         109,
				BuyUrl:        scraper.joinPath(scraper.baseUrl, "/buy/catId=123"),
			},
		}
		assert.Equal(t, want, got)
	})

	t.Run("Can find star hub event links", func(t *testing.T) {
		want := map[string]string{
			"AA": "/amon-amarth-tickets/performer/367989/",
			"BB": "/anirudh-tickets/category/135010106/",
			"CC": "/anvil-tickets/category/710749/",
		}
		got, err := scraper.FindStarHubEventLinks("/concert-tickets/category/1/")
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("return star hub all events ticket info", func(t *testing.T) {
		eventLink := "/red-hot-chili-peppers-tickets/performer/7527"
		want := []Event{
			{
				EventName:  "Red Hot Chili Peppers",
				Time:       "2023-02-19T15:00:00",
				Venue:      "Tokyo Dome;Tokyo, Japan",
				TicketLink: "/red-hot-chili-peppers-tokyo-tickets-2-19-2023/event/151207028/",
				Tickets: []Ticket{
					{
						Price:  229,
						BuyUrl: "/red-hot-chili-peppers-tokyo-tickets-2-19-2023/event/151207028/",
					},
				},
				Platform: "StarHub",
			},
			{
				EventName:  "Red Hot Chili Peppers",
				Time:       "2023-02-21T18:00:00",
				Venue:      "Osaka Jo Hall;Osaka, Japan",
				TicketLink: "/red-hot-chili-peppers-osaka-tickets-2-21-2023/event/151207029/",
				Tickets: []Ticket{
					{
						Price:  299,
						BuyUrl: "/red-hot-chili-peppers-osaka-tickets-2-21-2023/event/151207029/",
					},
				},
				Platform: "StarHub",
			},
		}
		got, err := scraper.GetStarHubEvents(eventLink)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
