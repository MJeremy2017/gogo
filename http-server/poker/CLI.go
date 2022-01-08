package poker

import (
	"io"
	"os"
	"bufio"
	"strings"
	"time"
	"fmt"
)


type CLI struct {
	store 	PlayerStore
	in 		*bufio.Scanner
	alerter BlindAlerter
}

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

// make this functional type implement BlindAlerter interface
type BlindAlerterFunc func(duration time.Duration, amount int)

func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int) {
	a(duration, amount)
}

// now this function can be a BlindAlerterFunc type which implements BlindAlerter interface
func StdOutAlerter(duration time.Duration, amount int) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(os.Stdout, "Blind is now %d\n", amount)
	})
}


func (c *CLI) PlayPoker() {
	c.scheduleBlindAlerts()
	input := c.readLine()

	winner := extractWinner(input)
	c.store.RecordWin(winner)
}

func (c *CLI) readLine() string {
	c.in.Scan()
	return c.in.Text()
}

func (c *CLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 
	1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second

	for _, blind := range blinds {
		c.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime += 10 * time.Minute
	}
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