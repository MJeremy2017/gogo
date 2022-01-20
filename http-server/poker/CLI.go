package poker

import (
	"io"
	"bufio"
	"strings"
	"time"
	"strconv"
	"fmt"
	"io/ioutil"
)

const PlayerPrompt = "Please enter the number of players: "
const BadPlayerInputErrMsg = "Bad value received for number of players, please try again with a number"
const BadWinnerInputErrMsg = "Bad value received for winner, please try again with {name} wins"
const Win = "wins"


type CLI struct {
	in 		*bufio.Scanner
	out 	io.Writer
	game	Game
}

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int, to io.Writer)
}

// make this functional type implement BlindAlerter interface
type BlindAlerterFunc func(duration time.Duration, amount int, to io.Writer)

func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	a(duration, amount, to)
}

// now this function can be a BlindAlerterFunc type which implements BlindAlerter interface
func Alerter(duration time.Duration, amount int, to io.Writer) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(to, "Blind is now %d\n", amount)
	})
}


func (c *CLI) PlayPoker() {
	fmt.Fprint(c.out, PlayerPrompt)

	numberOfPlayersInput := c.readLine()
	numberOfPlayers, err := strconv.Atoi(strings.Trim(numberOfPlayersInput, "\n"))
	if err != nil {
		fmt.Fprintf(c.out, BadPlayerInputErrMsg)
		return
	}
	fmt.Printf("number of players: %d\n", numberOfPlayers)

	c.game.Start(numberOfPlayers, ioutil.Discard)
	winnerInput := c.readLine()
	if !validateWinnerInput(winnerInput) {
		fmt.Fprintf(c.out, BadWinnerInputErrMsg)
		return
	}

	winner := extractWinner(winnerInput)
	c.game.Finish(winner)
}

func (c *CLI) readLine() string {
	c.in.Scan()
	return c.in.Text()
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func validateWinnerInput(winnerInput string) bool {
	tail := winnerInput[(len(winnerInput) - len(Win)):]
	if tail != Win {
		return false
	}
	return true
}

func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in: bufio.NewScanner(in),
		out: out,
		game: game,
	}
}