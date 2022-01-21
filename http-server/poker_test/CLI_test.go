package poker_test

import (
	"bytes"
	"fmt"
	"io"
	"server/poker"
	"strings"
	"testing"
	"time"
)

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

var dummyAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

func TestCLI(t *testing.T) {
	t.Run("record chris win from user input", func(t *testing.T) {
		in := userSends("5", "Chris wins")
		game := &poker.GameSpy{}
		stdOut := &bytes.Buffer{}

		cli := poker.NewCLI(in, stdOut, game)
		cli.PlayPoker()

		assertMessageSentToUser(t, stdOut, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 5)
		assertFinishCallWith(t, game, "Chris")
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := userSends("1", "Cleo wins")
		playerStore := &poker.StubPlayerStore{}

		game := poker.NewGame(dummyAlerter, playerStore)
		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Cleo")
	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := userSends("5", "Cleo wins")
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		game := poker.NewGame(blindAlerter, playerStore)
		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduleAlert(t, got, want)
			})
		}
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {

		stdOut := &bytes.Buffer{}
		in := strings.NewReader("desp\n")
		game := &poker.GameSpy{}

		cli := poker.NewCLI(in, stdOut, game)
		cli.PlayPoker()

		if game.StartCalled {
			t.Error("game shouldn't be started but get called")
		}

		assertMessageSentToUser(t, stdOut, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)

	})

	t.Run("it prompts user error and end the game when winner format is incorrect", func(t *testing.T) {
		in := userSends("5", "Chris is a killer")
		game := &poker.GameSpy{}
		stdOut := &bytes.Buffer{}

		cli := poker.NewCLI(in, stdOut, game)
		cli.PlayPoker()

		assertGameStartedWith(t, game, 5)
		assertMessageSentToUser(t, stdOut, poker.PlayerPrompt, poker.BadWinnerInputErrMsg)
	})

}

func assertScheduleAlert(t testing.TB, got, want scheduledAlert) {
	t.Helper()
	if got.amount != want.amount {
		t.Errorf("got amount %d, want %d", got.amount, want.amount)
	}

	if got.scheduledAt != want.scheduledAt {
		t.Errorf("got scheduled time %v, want %v", got.scheduledAt, want.scheduledAt)
	}
}

func assertMessageSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}

}

func assertGameStartedWith(t testing.TB, game *poker.GameSpy, numberOfPlayers int) {
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

func assertFinishCallWith(t testing.TB, game *poker.GameSpy, winner string) {
	t.Helper()
	want := winner
	got := game.FinishedWith
	if got != want {
		t.Errorf("wanted winner %s but got %s", want, got)
	}

}
