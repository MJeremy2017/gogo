package scrape

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
	mux.HandleFunc("/sg/Concert-Tickets", func(w http.ResponseWriter, r *http.Request) {
		data := `
<!doctype html>
<html>
<div class="bMag3 pbxxxl pb0-s">
        <div class="pgw">
            <h1 class="mbxs ibk t cWht xxl xl-s ">Concert Tickets</h1>
            <ul class="cloud">
                    <li><a href="/sg/Concert-Tickets/Clubs-and-Dance">Club and dance </a></li>
                    <li><a href="/sg/Theater-Tickets/Flamenco">Flamenco </a></li>
            </ul>
        </div>
</div>
</html>
`
		_, _ = fmt.Fprint(w, data)
	})

	mux.HandleFunc("/sg/Concert-Tickets/Clubs-and-Dance", func(w http.ResponseWriter, r *http.Request) {
		data := `
<html>
    <div class="uuxxl pgw">
		<a href="something">abc</a>
        <h2 class="txtc t xxl pbl">All Events</h2>
            <h3 class="t xxl"></h3>
            <ul class="cloud mbxl">
                    <li>
                        <a class="t" href="/sg/Concert-Tickets/Rock-and-Pop/5-Seconds-of-Summer-Tickets">5 Seconds of Summer</a>
                    </li>
            </ul>
            <h3 class="t xxl">B</h3>
            <ul class="cloud mbxl">
                    <li>
                        <a class="t" href="/sg/Concert-Tickets/Rock-and-Pop/Bastille-Tickets">Bastille</a>
                    </li>
                    <li>
                        <a class="t" href="/sg/Concert-Tickets/Rock-and-Pop/Brigitte-Tickets">Brigitte</a>
                    </li>
            </ul>
    </div>
</html>
`
		_, _ = fmt.Fprint(w, data)
	})

	return httptest.NewServer(mux)
}

func TestScraper_FindLinks(t *testing.T) {
	s := setUp()
	defer s.Close()
	scraper := NewScraper(s.URL)

	t.Run("Can find category links", func(t *testing.T) {
		want := map[string]string{
			"Concert Tickets":  "/sg/Concert-Tickets",
			"Sports Tickets":   "/sg/Sports-Tickets",
			"Theatre Tickets":  "/sg/Theatre-Tickets",
			"Festival Tickets": "/sg/Festival-Tickets",
		}
		got, err := scraper.FindLinks("", CategoryQuery)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Can find event type links", func(t *testing.T) {
		want := map[string]string{
			"Club and dance": "/sg/Concert-Tickets/Clubs-and-Dance",
			"Flamenco":       "/sg/Theater-Tickets/Flamenco",
		}
		got, err := scraper.FindLinks("/sg/Concert-Tickets", EventTypeQuery)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Can find events links", func(t *testing.T) {
		want := map[string]string{
			"5 Seconds of Summer": "/sg/Concert-Tickets/Rock-and-Pop/5-Seconds-of-Summer-Tickets",
			"Bastille":            "/sg/Concert-Tickets/Rock-and-Pop/Bastille-Tickets",
			"Brigitte":            "/sg/Concert-Tickets/Rock-and-Pop/Brigitte-Tickets",
		}
		got, err := scraper.FindLinks("/sg/Concert-Tickets/Clubs-and-Dance", EventQuery)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("return empty when no events", func(t *testing.T) {
		want := make(map[string]string)
		got, err := scraper.FindLinks("/", EventQuery)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	// TODO: got ticket of next layer
	t.Run("return all events ticket info", func(t *testing.T) {
		eventLink := "/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets"
		want := []Event{
			{
				EventName: "Super Junior Tickets",
			},
		}
		got, err := scraper.GetEvents(eventLink)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
