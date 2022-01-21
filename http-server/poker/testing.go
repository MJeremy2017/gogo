package poker

import (
	"fmt"
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

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

type scheduledAlert struct {
	scheduledAt time.Duration
	amount      int
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.scheduledAt)
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	s.alerts = append(s.alerts, scheduledAlert{duration, amount})
}

type GameSpy struct {
	StartedWith  int
	FinishedWith string
	StartCalled  bool
}

func (g *GameSpy) Start(numberOfPlayers int, alertsDestination io.Writer) {
	g.StartCalled = true
	g.StartedWith = numberOfPlayers
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
	want := numberOfPlayers
	got := game.StartedWith
	if got != want {
		t.Errorf("wanted start call with %d but got %d", want, got)
	}
}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}

func assertFinishCallWith(t testing.TB, game *GameSpy, winner string) {
	t.Helper()
	want := winner
	got := game.FinishedWith
	if got != want {
		t.Errorf("wanted winner %s but got %s", want, got)
	}

}
