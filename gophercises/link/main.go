package main

import (
	"flag"
	"fmt"
	"link/parser"
)

// TODO add cmd flags
// TODO add xml parser

func main() {
	startUrl := flag.String("source", "", "starting crawling url")
	maxDepth := flag.Int("depth", 2, "maximum crawling depth")
	flag.Parse()

	got := parser.BrowseLinks(*startUrl, *maxDepth)
	fmt.Println("got", got)
}
