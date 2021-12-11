package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"fmt"
)


type StubPlayerStore struct {
	scores map[string]int
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	return s.scores[name]
}


func TestGETPlayers(t *testing.T) {
	stubPlayerScore := &StubPlayerStore{
		scores: map[string]int{
			"Pepper": 20,
			"Floyd": 10,
		},
	}
	server := &PlayerServer{stubPlayerScore}

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "20"

		assertBodyEqual(t, got, want)
	})
	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "10"
		assertBodyEqual(t, got, want)
	})



}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertBodyEqual(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q wnat %q", got, want)
	}
}