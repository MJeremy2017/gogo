package race

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"time"
)

func TestRacer(t *testing.T) {
	t.Run("Test return faster server", func(t *testing.T) {
		slowServer := makeDelayedServer(20 * time.Millisecond)
		fastServer := makeDelayedServer(0)

		// call function at the end of the scope
		defer slowServer.Close()
		defer fastServer.Close()

		slowURL := slowServer.URL
	    fastURL := fastServer.URL

	    want := fastURL
	    got, err := Racer(slowURL, fastURL)
	    if err != nil {
	    	t.Fatalf("did not expect an error but got one %v", err)
	    }
	    
	    if got != want {
	    	t.Errorf("got %q want %q", got, want)
	    }
	})

	t.Run("Test return error after a specified time", func(t *testing.T) {
		slowServer := makeDelayedServer(20 * time.Millisecond)
		fastServer := makeDelayedServer(30 * time.Millisecond)
		timeout := 10 * time.Millisecond

		defer slowServer.Close()
		defer fastServer.Close()

		slowURL := slowServer.URL
	    fastURL := fastServer.URL

	    _, err := ConfigurableRacer(slowURL, fastURL, timeout)
	    if err == nil {
	    	t.Errorf("expected an timeout error but got nil")
	    }
	})


}

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
			time.Sleep(delay)
			w.WriteHeader(http.StatusOK)
		}))
}