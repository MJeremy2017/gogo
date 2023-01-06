package scrape

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScraper_FindCategory(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := `
<ul class="prinav">
        <li class="t cat3">
            <a href="/sg/Concert-Tickets">Concert Tickets</a>
        </li>
        <li class="t cat2">
            <a href="/sg/Sports-Tickets">Sports Tickets</a>
        </li>
        <li class="t cat1">
            <a href="/sg/Theater-Tickets">Theatre Tickets</a>
        </li>
        <li class="t cat1023">
            <a href="/sg/Festival-Tickets">Festival Tickets</a>
        </li>
</ul>
`
		_, _ = fmt.Fprint(w, data)
	})

	s := httptest.NewServer(mux)
	defer s.Close()

	scraper := &Scraper{}
	want := map[string]string{
		"Concert Tickets":  "/sg/Concert-Tickets",
		"Sports Tickets":   "/sg/Sports-Tickets",
		"Theatre Tickets":  "/sg/Theatre-Tickets",
		"Festival Tickets": "/sg/Festival-Tickets",
	}
	got := scraper.FindCategory()

	assert.Equal(t, want, got)
}
