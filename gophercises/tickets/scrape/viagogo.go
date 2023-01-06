package scrape

import (
	"net/url"
)

// scrape all the categories
const baseUrl = "https://www.viagogo.com"

type Scraper struct {
}

func (s *Scraper) FindCategory() map[string]url.URL {
	//_ := colly.NewCollector()

	return nil
}
