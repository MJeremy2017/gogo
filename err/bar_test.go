package err

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"fmt"
	"errors"
)


type BadStatusError struct {
	URL 	string
	Status 	int
}


func (b BadStatusError) Error() string {
	return fmt.Sprintf("did not get 200 from %s, got %d", b.URL, b.Status)
}


func TestDumbGetter(t *testing.T) {
	t.Run("get a status error", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	        res.WriteHeader(http.StatusTeapot)
	    }))
	    defer svr.Close()

	    _, err := DumbGetter(svr.URL)

	    if err == nil {
	        t.Fatal("expected an error")
	    }

	    var got BadStatusError
	    isStatusError := errors.As(err, &got)
	    if !isStatusError {
	    	t.Fatalf("was not a BadStatusError")
	    }

	    want := BadStatusError{
	    	svr.URL,
	    	http.StatusTeapot,
	    }

	    if got != want {
	        t.Errorf("got %v, want %v", got, want)
	    }
	})
	
}