package link

import (
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
				trimmedData := strings.TrimSpace(n.Data)
				if n.Type == html.TextNode && trimmedData != "" {
					innerText += trimmedData + " "
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
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			dfs(c)
		}
	}

	dfs(node)
	return links, nil
}
