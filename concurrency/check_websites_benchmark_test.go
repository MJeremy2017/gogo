package concurrency

import (
	"testing"
	"time"
)

func slowStubWebsiteChecker(_ string) bool {
	time.Sleep(20 * time.Millisecond)
	return true
}

func BenchmarkCheckWebsites(b *testing.B) {
	urls := make([]string, 100)
	for i := 0; i < 100; i++ {
		urls[i] = "a url"
	}

	for i := 0; i < b.N; i++ {
		// only inside is the function to be tested
		CheckWebsites(slowStubWebsiteChecker, urls)
	}
}