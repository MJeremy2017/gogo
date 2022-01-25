package poker

import (
	"io"
	"strings"
	"testing"
	"time"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	return s.scores[name]
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

type GameSpy struct {
	StartCalled  bool
	StartedWith  int
	BlindAlert   []byte

	FinishedCalled	bool
	FinishedWith 	string
}

func (g *GameSpy) Start(numberOfPlayers int, alertsDestination io.Writer) {
	g.StartCalled = true
	g.StartedWith = numberOfPlayers
	alertsDestination.Write(g.BlindAlert)
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
}

// AssertPlayerWin allows you to spy on the store's calls to RecordWin.
func AssertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != winner {
		t.Errorf("did not store the correct winner got %q want %q", store.winCalls[0], winner)
	}
}

func assertGameStartedWith(t testing.TB, game *GameSpy, numberOfPlayers int) {
	t.Helper()

	passed := retryUntil(500 * time.Millisecond, func() bool {
		return numberOfPlayers == game.StartedWith
	})

	want := numberOfPlayers
	got := game.StartedWith
	if !passed {
		t.Errorf("wanted start call with %d but got %d", want, got)
	}
}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}

func assertFinishCallWith(t testing.TB, game *GameSpy, winner string) {
	t.Helper()

	passed := retryUntil(500 * time.Millisecond, func() bool {
		return game.FinishedWith == winner
	})

	if !passed {
		t.Errorf("wanted winner %s but got %s", winner, game.FinishedWith)
	}

}

func within(t testing.TB, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <- time.After(d):
		t.Error("time out")
	case <- done:
	}
}


func retryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}





