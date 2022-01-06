package poker_test

import (
	"server/poker"
	"strings"
	"testing"
)


func TestCLI(t *testing.T) {

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}
		dummyAlerter := &poker.SpyBlindAlerter{}

		cli := poker.NewCLI(playerStore, in, dummyAlerter)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &poker.StubPlayerStore{}
		dummyAlerter := &poker.SpyBlindAlerter{}

		cli := poker.NewCLI(playerStore, in, dummyAlerter)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Cleo")
	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &poker.SpyBlindAlerter{}

		cli := poker.NewCLI(playerStore, in, blindAlerter)
		cli.PlayPoker()

		if len(blindAlerter.Alerts) != 1 {
			t.Fatal("expected a blind alert but got none")
		}

	})

}
