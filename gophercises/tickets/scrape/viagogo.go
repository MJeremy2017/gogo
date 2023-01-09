package scrape

import (
	"github.com/gocolly/colly"
	"strings"
)

const categoryQuery = ".prinav a[href]"
const eventTypeQuery = ".cloud a[href]"

type Scraper struct {
	baseUrl string
}

func NewScraper(baseUrl string) *Scraper {
	return &Scraper{baseUrl: baseUrl}
}

// FindCategory returns a map with key `category name` and value `relative path after host:port`
func (s *Scraper) FindCategory() (map[string]string, error) {
	res := make(map[string]string)

	c := colly.NewCollector()
	c.OnHTML(categoryQuery, func(e *colly.HTMLElement) {
		relativePath := e.Attr("href")
		text := strings.TrimSpace(e.Text)
		res[text] = relativePath
	})

	err := c.Visit(s.baseUrl)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Scraper) FindEventType() (map[string]string, error) {
	res := make(map[string]string)

	c := colly.NewCollector()
	c.OnHTML(eventTypeQuery, func(e *colly.HTMLElement) {
		relativePath := e.Attr("href")
		text := strings.TrimSpace(e.Text)
		res[text] = relativePath
	})

	err := c.Visit(s.baseUrl)
	if err != nil {
		return nil, err
	}
	return res, nil
}
