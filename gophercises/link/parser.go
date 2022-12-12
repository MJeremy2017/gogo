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

	var links []Link
	var dfs func(n *html.Node)

	dfs = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			lk := p.parseLinkFromNode(n)
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

func (p *Parser) GetLinkFromAttr(attr []html.Attribute) string {
	for _, a := range attr {
		if a.Key == "href" {
			return a.Val
		}
	}
	return ""
}

func (p *Parser) parseLinkFromNode(n *html.Node) Link {
	var innerText string
	var strDfs func(n *html.Node)

	strDfs = func(n *html.Node) {
		if n.Type == html.TextNode {
			if !strings.HasSuffix(innerText, " ") {
				innerText += " "
			}
			innerText += strings.TrimSpace(n.Data)
		}
		for cc := n.FirstChild; cc != nil; cc = cc.NextSibling {
			strDfs(cc)
		}
	}
	strDfs(n)
	return Link{
		Href: p.GetLinkFromAttr(n.Attr),
		Text: strings.TrimSpace(innerText),
	}
}
