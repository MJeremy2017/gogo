package scrape_test

import (
	"fmt"
	"github.com/gocolly/colly"
	"testing"
)

// TODO start from concert tickets and see how to find categories
func TestA(t *testing.T) {
	depthFunc := colly.MaxDepth(2)
	c := colly.NewCollector(depthFunc)

	// Find and visit all links
	c.OnHTML(".prinav a[href]", func(e *colly.HTMLElement) {
		fmt.Printf("attr %s text %s\n", e.Attr("href"), e.Text)
		fmt.Println("url", e.Request.URL.Path)
		e.Request.Visit(e.Attr("href"))
	})

	c.Visit("https://www.viagogo.com/sg")
}