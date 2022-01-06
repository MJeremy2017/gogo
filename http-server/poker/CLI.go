package poker

import (
	"io"
	"bufio"
	"strings"
	"time"
)


type CLI struct {
	store 	PlayerStore
	in 		*bufio.Scanner
	alerter BlindAlerter
}

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

func (c *CLI) PlayPoker() {
	input := c.readLine()

	winner := extractWinner(input)
	c.store.RecordWin(winner)
}

func (c *CLI) readLine() string {
	c.in.Scan()
	return c.in.Text()
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func NewCLI(store PlayerStore, in io.Reader, alerter BlindAlerter) *CLI {
	return &CLI{
		store: store,
		in: bufio.NewScanner(in),
		alerter: alerter,
	}
}