package link

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
)

type Parser struct {
	html io.Reader
}

type Link struct {
	Href string
	Text string
}

func NewParser(html io.Reader) *Parser {
	return &Parser{html}
}

func (p *Parser) ParseLinks() ([]Link, error) {
	// TODO parse element fix basic test
	node, err := html.Parse(p.html)
	if err != nil {
		return nil, err
	}

	var dfs func(n *html.Node)
	var links []Link
	dfs = func(n *html.Node) {
		fmt.Printf("type: %+v, data: %+v, atom: %+v, attr: %+v", n.Type, n.Data, n.DataAtom.String(), n.Attr)
		if n.Type == html.ElementNode && n.Data == "a" {
			lk := Link{
				Href: n.Attr[0].Val,
				Text: "",
			}
			links = append(links, lk)
		}
		c := n.FirstChild
		for c != nil {
			dfs(c)
			c = c.NextSibling
		}
	}

	dfs(node)
	return links, nil
}
