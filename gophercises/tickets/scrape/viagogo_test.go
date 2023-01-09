package scrape

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
)

func setUp() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := `
<!doctype html>
<html>
<ul class="prinav">
        <li class="t cat3">
            <a href="/sg/Concert-Tickets">Concert Tickets</a>
        </li>
        <li class="t cat2">
            <a href="/sg/Sports-Tickets">Sports Tickets</a>
        </li>
        <li class="t cat1">
            <a href="/sg/Theatre-Tickets">Theatre Tickets</a>
        </li>
        <li class="t cat1023">
            <a href="/sg/Festival-Tickets">Festival Tickets</a>
        </li>
</ul>
</html>
`
		_, _ = fmt.Fprint(w, data)
	})

	return httptest.NewServer(mux)
}

func TestScraper_FindCategory(t *testing.T) {
	s := setUp()
	defer s.Close()

	scraper := NewScraper(s.URL)

	got := map[string]string{
		"Concert Tickets":  path.Join(s.URL, "/sg/Concert-Tickets"),
		"Sports Tickets":   path.Join(s.URL, "/sg/Sports-Tickets"),
		"Theatre Tickets":  path.Join(s.URL, "/sg/Theatre-Tickets"),
		"Festival Tickets": path.Join(s.URL, "/sg/Festival-Tickets"),
	}
	want, err := scraper.FindCategory()
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
