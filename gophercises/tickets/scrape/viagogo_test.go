package scrape_test

import (
	"fmt"
	"github.com/gocolly/colly"
	"testing"
)

func TestA(t *testing.T) {
	fmt.Println("h")
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://go-colly.org/")
}
