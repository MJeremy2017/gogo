package main

import (
	"io"
	"fmt"
	"os"
	"time"
)

const finalWord = "Go!"
const countDownStart = 3


type Sleeper interface {
	Sleep()
}

type DefaultSleeper struct {}

type ConfigurableSleeper struct {
	duration time.Duration
	sleep func(time.Duration)
}

func (d DefaultSleeper) Sleep() {
	time.Sleep(1 * time.Second)
}

func (c *ConfigurableSleeper) Sleep() {
	c.sleep(c.duration)
}


func Countdown(writer io.Writer, sleeper Sleeper) {
	for i := countDownStart; i > 0; i-- {
		sleeper.Sleep()
		fmt.Fprintln(writer, i)
	}
	sleeper.Sleep()
	fmt.Fprint(writer, finalWord)
}

func main() {
	sleeper := &ConfigurableSleeper{time.Second, time.Sleep}
	Countdown(os.Stdout, sleeper)
}