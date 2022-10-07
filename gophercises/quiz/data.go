package main

import "fmt"

type GameTally struct {
	correct int32
	total   int32
}

func (g *GameTally) increaseCorrect() {
	g.correct += 1
}

func (g *GameTally) increaseTotal() {
	g.total += 1
}

func (g *GameTally) printScore() {
	fmt.Printf("You have got %d out of %d quiz right\n", g.correct, g.total)
}

func NewGameTally() *GameTally {
	return &GameTally{0, 0}
}
