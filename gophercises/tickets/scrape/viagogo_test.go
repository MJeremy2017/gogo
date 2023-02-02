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
		got, err := scraper.FindLinks("", CategoryQuery)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Can find event type links", func(t *testing.T) {
		want := map[string]string{
			"Club and dance": "/sg/Concert-Tickets/Clubs-and-Dance",
			"Flamenco":       "/sg/Theater-Tickets/Flamenco",
		}
		got, err := scraper.FindLinks("/sg/Concert-Tickets", EventTypeQuery)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Can find events links", func(t *testing.T) {
		want := map[string]string{
			"5 Seconds of Summer": "/sg/Concert-Tickets/Rock-and-Pop/5-Seconds-of-Summer-Tickets",
			"Bastille":            "/sg/Concert-Tickets/Rock-and-Pop/Bastille-Tickets",
			"Brigitte":            "/sg/Concert-Tickets/Rock-and-Pop/Brigitte-Tickets",
		}
		got, err := scraper.FindLinks("/sg/Concert-Tickets/Clubs-and-Dance", EventQuery)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("return empty when no events", func(t *testing.T) {
		want := make(map[string]string)
		got, err := scraper.FindLinks("/", EventQuery)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	// TODO: get ticket price for different quantities https://www.viagogo.com/sg/Concert-Tickets/Rock-and-Pop/Mayday-Tickets/E-151357967?qty=1,2...
	t.Run("return all events ticket info", func(t *testing.T) {
		eventLink := "/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets"
		want := []Event{
			{
				EventName:  "Super Junior Tickets",
				Time:       "2023-02-09T20:00:00",
				Venue:      "Sao Paulo, Brazil",
				TicketLink: "/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets/E-151336327",
			},
			{
				EventName:  "Super Junior Tickets",
				Time:       "2023-02-18T19:30:00",
				Venue:      "Selangor, Malaysia",
				TicketLink: "/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets/E-151396140",
			},
		}
		got, err := scraper.GetEvents(eventLink)
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
			},
		}
		assert.Equal(t, want, got)
	})
}
