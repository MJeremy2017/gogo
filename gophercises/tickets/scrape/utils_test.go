package scrape

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScraper_SortEvents(t *testing.T) {
	t.Run("can sort all events by ticket price", func(t *testing.T) {
		events := []Event{
			{
				EventName:  "Super Junior Tickets",
				Time:       "2023-02-09T20:00:00",
				Venue:      "Sao Paulo, Brazil",
				TicketLink: "/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets/E-151336327",
				Tickets: []Ticket{
					{
						QuantityRange: "1 - 4",
						Price:         109,
					},
					{
						QuantityRange: "1 - 4",
						Price:         124,
					},
					{
						QuantityRange: "1 - 4",
						Price:         118,
					},
				},
			},
		}

		want := []Event{
			{
				EventName:  "Super Junior Tickets",
				Time:       "2023-02-09T20:00:00",
				Venue:      "Sao Paulo, Brazil",
				TicketLink: "/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets/E-151336327",
				Tickets: []Ticket{
					{
						QuantityRange: "1 - 4",
						Price:         109,
					},
					{
						QuantityRange: "1 - 4",
						Price:         118,
					},
					{
						QuantityRange: "1 - 4",
						Price:         124,
					},
				},
			},
		}

		got := SortEventTicketsByPrice(events)

		assert.Equal(t, got, want)

	})
}
