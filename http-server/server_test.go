package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"fmt"
	"reflect"
	"encoding/json"
	"io"
)


type StubPlayerStore struct {
	scores map[string]int
	winCalls []string
	league []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	return s.scores[name]
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() []Player {
	return s.league
}


func TestGETPlayers(t *testing.T) {
	stubPlayerScore := &StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd": 10,
		},
		nil,
		nil,
	}
	server := NewPlayerServer(stubPlayerScore)

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		got := response.Body.String()
		want := "20"

		assertBodyEqual(t, got, want)
	})
	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		got := response.Body.String()
		want := "10"
		assertBodyEqual(t, got, want)
	})
	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})

}

func TestStoreWins(t *testing.T) {
	stub := &StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := NewPlayerServer(stub)
	
	t.Run("records when post", func(t *testing.T) {
		player := "Pepper"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(stub.winCalls) != 1 {
			t.Errorf("got %d calls to Recording Win want %d", len(stub.winCalls), 1)
		}

		if stub.winCalls[0] != player {
			t.Errorf("did not store correct winner, got %q want %q", stub.winCalls[0], player)
		}
	})
}

func TestRecordingWinsAndRetrivingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))
	assertStatus(t, response.Code, http.StatusOK)

	assertBodyEqual(t, response.Body.String(), "3")
}


func TestLeague(t *testing.T) {
	t.Run("returns json on /league", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 22},
			{"Chris", 13},
			{"Neo", 31},
		}

		stub := &StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(stub)

		request := newGetLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		assertStatus(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)
	})
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)
	if err != nil {
		t.Fatalf("Unable to parse response body %q to player, '%v'", body, err)
	}
	return
}


func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got status %d want %d", got, want)
	}

}

func assertLeague(t testing.TB, got, want []Player) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}


func newGetLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/league"), nil)
	return req
}


func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertBodyEqual(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}