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

	// TODO: add multiple links cases to testing
	t.Run("return all events ticket info", func(t *testing.T) {
		eventLink := "/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets"
		want := []Event{
			{
				EventName:  "Super Junior Tickets",
				TicketLink: "/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets/E-151336327",
			},
			{
				EventName:  "Super Junior Tickets",
				TicketLink: "/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets/E-151396140",
			},
		}
		got, err := scraper.GetEvents(eventLink)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
