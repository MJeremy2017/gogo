package main

import (
	"fmt"
	"link/parser"
)

func main() {
	u := "https://www.google.com"
	got := parser.BrowseLinks(u, 2)
	fmt.Println("got", got)
}
