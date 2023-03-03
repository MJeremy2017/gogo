package scrape_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"tickets/scrape"
)

func TestViaGogo(t *testing.T) {
	events, err := scrape.LoadJsonToEvents("viagogo_event.json")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(events[0])
}

func TestStarHub(t *testing.T) {
	events, err := scrape.LoadJsonToEvents("stubhub_event.json")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(events[0])
}

func TestAsync(t *testing.T) {
	// TODO
	//var (
	//	mu     sync.Mutex
	//	number int
	//)
	//
	//func main() {
	//	// start the side calculation in a separate goroutine
	//	go func() {
	//		for {
	//			// do the side calculation
	//			time.Sleep(5 * time.Second)
	//			mu.Lock()
	//			number++
	//			mu.Unlock()
	//		}
	//	}()
	//
	//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//		// acquire the lock to access the shared state
	//		mu.Lock()
	//		defer mu.Unlock()
	//
	//		// display the calculated number on the front page
	//		tmpl := template.Must(template.ParseFiles("index.html"))
	//		tmpl.Execute(w, number)
	//	})
	//
	//	fmt.Println("Starting server on http://localhost:8080")
	//	http.ListenAndServe(":8080", nil)
	//}
}

func getAndSaveResponse(url string) {
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		log.Println(err)
	}

	bytes, _ := ioutil.ReadAll(response.Body)

	f, _ := os.Create("test.html")
	defer f.Close()
	f.WriteString(string(bytes))
	fmt.Println(string(bytes))
}

func postAndSaveResponse(url string) {
	response, err := http.Post(url, "application/json", nil)
	defer response.Body.Close()
	if err != nil {
		log.Println(err)
	}

	bytes, _ := ioutil.ReadAll(response.Body)

	f, _ := os.Create("test.json")
	defer f.Close()
	f.WriteString(string(bytes))
	fmt.Println(string(bytes))
}
