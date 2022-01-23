package poker

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

var dummyGame = &GameSpy{}

func TestGETPlayers(t *testing.T) {
	stubPlayerScore := &StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := mustMakePlayerServer(t, stubPlayerScore, dummyGame)

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
	server := mustMakePlayerServer(t, stub, dummyGame)

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
			t.Errorf("did not store the correct winner, got %q want %q", stub.winCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {
	t.Run("returns json on /league", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 22},
			{"Chris", 13},
			{"Neo", 31},
		}

		stub := &StubPlayerStore{nil, nil, wantedLeague}
		server := mustMakePlayerServer(t, stub, dummyGame)

		request := newGetLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		// assert header
		assertContentType(t, response, jsonContentType)

		got := getLeagueFromResponse(t, response.Body)
		assertStatus(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)
	})
}

func TestGame(t *testing.T) {
	t.Run("Get /game returns 200", func(t *testing.T) {
		stub := &StubPlayerStore{}
		server := mustMakePlayerServer(t, stub, dummyGame)

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("start a playGame with 3 players and declare Ruth the winner", func(t *testing.T) {
		wantedBlindAlert := "Blind is 100"
		winner := "Ruth"

		game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}
		store := &StubPlayerStore{}
		playerServer := mustMakePlayerServer(t, store, game)
		server := httptest.NewServer(playerServer)
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
		ws := mustDialWS(t, wsURL)
		defer ws.Close()

		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		time.Sleep(10 * time.Millisecond)
		assertGameStartedWith(t, dummyGame, 3)
		assertFinishCallWith(t, dummyGame, winner)

		_, gotBlindAlert, _ := ws.ReadMessage()
		if wantedBlindAlert != string(gotBlindAlert) {
			t.Errorf("got blind alert %q, want %q", gotBlindAlert, wantedBlindAlert)
		}

	})

}

func mustMakePlayerServer(t *testing.T, store PlayerStore, game *GameSpy) *PlayerServer {
	server, err := NewPlayerServer(store, game)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}
	return ws
}

func writeWSMessage(t testing.TB, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Player) {
	t.Helper()
	league, err := NewLeague(body)
	if err != nil {
		t.Fatalf("Unable to parse response body %q to player, '%v'", body, err)
	}
	return
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type application/json, got %v",
			response.Result().Header)
	}
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

func newGameRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/game"), nil)
	return req
}

func assertBodyEqual(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
