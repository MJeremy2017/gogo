package link

import (
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

func (p *Parser) ParseLinks() []Link {
	_, _ = html.Parse(p.html)
	return nil
}
