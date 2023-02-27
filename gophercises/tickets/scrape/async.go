package scrape

var events []Event

func AsyncScrapeSiteEvents(host, fp string) {
	// TODO infinite loop
	chanEvents := make(chan []Event, 1)
	go func() {
		s := NewScraper(host)
		// TODO fix the if else
		if host == "https://www.stubhub.com" {
			events = s.GetStarHubAllEvents()
		} else {
			events = s.GetViagogoAllEvents()
		}
		chanEvents <- events
	}()
	events = <-chanEvents
	SaveEventsToJson(events, fp)
}
