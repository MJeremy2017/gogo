package mycontext

import (
	"net/http/httptest"
	"net/http"
	"testing"
	"time"
	"context"
)

type SpyStore struct {
	response string
	cancelled bool
}

func (s *SpyStore) Fetch() string {
	time.Sleep(20 * time.Millisecond)
	return s.response
}

func (s *SpyStore) Cancel() {
	s.cancelled = true
}


func TestServer(t *testing.T) {
	t.Run("Server can run", func(t *testing.T) {
		data := "Hello, World"
		stub := SpyStore{response: data}

		svr := Server(&stub)
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		if response.Body.String() != data {
			t.Errorf("got %q, want %q", response.Body.String(), data)
		}

		if stub.cancelled {
			t.Error("it should not have cancelled the store")
		}
	}) 

	t.Run("Tells store to cancel work if request is cancelled", func(t *testing.T) {
		data := "Hello, World"
		stub := SpyStore{response: data}

		svr := Server(&stub)
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		cancellingCtx, cancelFunc := context.WithCancel(request.Context())
		time.AfterFunc(5 * time.Millisecond, cancelFunc)
		request = request.WithContext(cancellingCtx)

		response := httptest.NewRecorder()
		svr.ServeHTTP(response, request)

		if !stub.cancelled {
			t.Error("expected request cancelled but not")
		}

	})

}