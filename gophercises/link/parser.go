package link

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"strings"
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
	node, err := html.Parse(p.html)
	if err != nil {
		return nil, err
	}

	var dfs func(n *html.Node)
	var links []Link
	dfs = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			var innerText string
			var strDfs func(n *html.Node)
			strDfs = func(n *html.Node) {
				fmt.Printf("inside %+v, %v", n.Type, n.Data)
				if n.Type == html.TextNode {
					innerText += strings.TrimSpace(n.Data) + " "
				}
				for cc := n.FirstChild; cc != nil; cc = cc.NextSibling {
					strDfs(cc)
				}
			}
			strDfs(n)
			lk := Link{
				Href: n.Attr[0].Val,
				Text: strings.TrimSpace(innerText),
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
