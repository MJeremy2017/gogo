package scrape

import (
	"github.com/gocolly/colly"
	"path"
)

const categoryQuery = ".prinav a[href]"

type Scraper struct {
	baseUrl string
}

func NewScraper(baseUrl string) *Scraper {
	return &Scraper{baseUrl: baseUrl}
}

// FindCategory returns a map with key `category name` and value `full path of url`
func (s *Scraper) FindCategory() (map[string]string, error) {
	res := make(map[string]string)

	c := colly.NewCollector()
	c.OnHTML(categoryQuery, func(e *colly.HTMLElement) {
		relativePath := e.Attr("href")
		res[e.Text] = path.Join(e.Request.URL.String(), relativePath)
	})

	err := c.Visit(s.baseUrl)
	if err != nil {
		return nil, err
	}
	return res, nil
}
