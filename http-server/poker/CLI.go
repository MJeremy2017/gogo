package poker

import (
	"io"
	"bufio"
	"strings"
)


type CLI struct {
	store 	PlayerStore
	in 		*bufio.Scanner
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

func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{
		store: store,
		in: bufio.NewScanner(in),
	}
}