package main

import (
	"flag"
	"link/parser"
	"log"
	"os"
)

func main() {
	startUrl := flag.String("source", "", "starting crawling url")
	maxDepth := flag.Int("depth", 2, "maximum crawling depth")
	flag.Parse()

	links := parser.BrowseLinks(*startUrl, *maxDepth)
	out, err := os.OpenFile("./links.xml", os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}
	err = parser.EncodeLinksToXML(links, out)
	if err != nil {
		log.Fatal(err)
	}
}
