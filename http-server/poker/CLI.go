package poker

import (
	"io"
	"os"
	"bufio"
	"strings"
	"time"
	"strconv"
	"fmt"
)

const PlayerPrompt = "Please enter the number of players: "
const BadPlayerInputErrMsg = "Bad value received for number of players, please try again with a number"


type CLI struct {
	in 		*bufio.Scanner
	out 	io.Writer
	game	Game
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
	fmt.Fprint(c.out, PlayerPrompt)

	numberOfPlayersInput := c.readLine()
	numberOfPlayers, err := strconv.Atoi(strings.Trim(numberOfPlayersInput, "\n"))
	if err != nil {
		fmt.Fprintf(c.out, BadPlayerInputErrMsg)
		return
	}
	fmt.Printf("number of players: %d\n", numberOfPlayers)

	c.game.Start(numberOfPlayers)
	winnerInput := c.readLine()

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

func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in: bufio.NewScanner(in),
		out: out,
		game: game,
	}
}