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
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 
	1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second

	for _, blind := range blinds {
		c.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime += 10 * time.Minute
	}

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