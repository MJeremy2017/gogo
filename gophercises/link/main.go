package main

import (
	"fmt"
	"link/parser"
)

func main() {
	url := "https://www.google.com"
	got := parser.BrowseLinks(url)
	fmt.Println("got", got)
}
