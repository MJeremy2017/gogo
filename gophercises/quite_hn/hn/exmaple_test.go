// The client tests need to be inside the hn package so this test file
// was created for examples that use the hn package from an external package
// like a normal user would.
package hn_test

import (
	"fmt"
	"hn/hn"
	"log"
	"math"
	"testing"
	"time"
)

const numStories = 30

// Test public facing interface
func TestAsyncClient(t *testing.T) {
	var client hn.Client

	startTime := time.Now()
	stories, err := client.GetTopStories(numStories)
	if err != nil {
		log.Fatalln(err)
	}
	for _, s := range stories {
		fmt.Printf("%s (by %s)\n", s.Title, s.By)
	}
	fmt.Printf("Total time %vs\n", NanoToSeconds(time.Since(startTime)))
}

func TestCachedClient(t *testing.T) {
	client := hn.NewCacheClient(60)

	startTime := time.Now()
	stories, err := client.GetTopStories(numStories)
	if err != nil {
		log.Fatalln(err)
	}
	for _, s := range stories {
		fmt.Printf("%s (by %s)\n", s.Title, s.By)
	}
	fmt.Printf("Total time %vs\n", NanoToSeconds(time.Since(startTime)))
}

func TestClient(t *testing.T) {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		panic(err)
	}

	startTime := time.Now()
	for i := 0; i < 5; i++ {
		item, err := client.GetItem(ids[i])
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s (by %s)\n", item.Title, item.By)
	}
	fmt.Printf("Total time %vs\n", NanoToSeconds(time.Since(startTime)))
}

func NanoToSeconds(t time.Duration) float64 {
	return float64(t) / math.Pow(10, 9)
}
