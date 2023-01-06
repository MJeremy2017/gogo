package scrape_test

import (
	"fmt"
	"github.com/gocolly/colly"
	"testing"
)

// TODO start from concert tickets and see how to find categories
// TODO in the main page, find all classname = "prinav"
func TestA(t *testing.T) {
	fmt.Println("h")
	depthFunc := colly.MaxDepth(2)
	c := colly.NewCollector(depthFunc)

	// Find and visit all links
	c.OnHTML(".prinav a[href]", func(e *colly.HTMLElement) {
		fmt.Printf("attr %s text %s\n", e.Attr("href"), e.Text)
		fmt.Println("url", e.Request.URL)
		e.Request.Visit(e.Attr("href"))
	})

	//c.OnRequest(func(r *colly.Request) {
	//	fmt.Println("Visiting", r.URL)
	//})

	c.Visit("https://www.viagogo.com/sg")
}
