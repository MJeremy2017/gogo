package race

import (
	"time"
	"net/http"
	"fmt"
)

func Racer(url1 string, url2 string, timeout time.Duration) (string, error) {
	// select wait on multiple channels and the first one wins
	select {
	case <-ping(url1):
		return url1, nil
	case <-ping(url2):
		return url2, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("time out waiting for %s and %s", url1, url2)
	}

}


func Racer2(url1 string, url2 string) string {
	ch := make(chan string)
	go func() {
		http.Get(url1)
		ch <- url1
	}()

	go func() {
		http.Get(url2)
		ch <- url2
	}()

	res := <-ch
	return res
}


func MeasureReponseTime(url string) time.Duration {
	start := time.Now()
	http.Get(url)
	duration := time.Since(start)
	return duration
}

func ping(url string) chan struct{} {
	ch := make(chan struct{})
	go func() {
		http.Get(url)
		close(ch)
	}()

	return ch
}