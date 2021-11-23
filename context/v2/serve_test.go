package mycontext

import (
	"net/http/httptest"
	"net/http"
	"testing"
	"time"
	"context"
	"errors"
)

type SpyStore struct {
	response string
	t testing.TB
}

type SpyRespnseWriter struct {
	written bool
}

func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
	data := make(chan string, 1)

	go func() {
		var result string
		// simulate incrementally adding result
		for _, c := range s.response {
			select {
			case <- ctx.Done():
				s.t.Log("spy store got cancelled")
				return
			default:
				time.Sleep(10 * time.Millisecond)
				result += string(c)
			}
		}
		data <- result
	}()


	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-data:
		return res, nil
	}

}

func (s *SpyRespnseWriter) Header() http.Header {
	s.written = true
	return nil
}

func (s *SpyRespnseWriter) Write([]byte) (int, error) {
	s.written = true
	return 0, errors.New("not implemented")
}

func (s *SpyRespnseWriter) WriteHeader(statusCode int) {
	s.written = true
}



func TestServer(t *testing.T) {
	t.Run("Server can run", func(t *testing.T) {
		data := "Hello, World"
		stub := &SpyStore{response: data, t: t}

		svr := Server(stub)
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		if response.Body.String() != data {
			t.Errorf("got %q, want %q", response.Body.String(), data)
		}

	}) 

	t.Run("Tells store to cancel work if request is cancelled", func(t *testing.T) {
		data := "Hello, World"
		stub := SpyStore{response: data, t: t}

		svr := Server(&stub)
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		cancellingCtx, cancelFunc := context.WithCancel(request.Context())
		time.AfterFunc(5 * time.Millisecond, cancelFunc)
		// add cancelling context to request
		request = request.WithContext(cancellingCtx)

		response := &SpyRespnseWriter{}
		svr.ServeHTTP(response, request)

		if response.written {
			t.Error("a response should not be written")
		}

	})

}

