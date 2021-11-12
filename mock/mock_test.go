package main

import (
	"testing"
	"bytes"
	"reflect"
	"time"
)

const sleep = "sleep"
const write = "write"

// Implement both Writer and Sleep Inferface
type CountdownOperationSpy struct {
	Calls []string
}

type SpyTime struct {
	duration time.Duration
}

func (s *SpyTime) Sleep(duration time.Duration) {
	s.duration = duration
}

func (c *CountdownOperationSpy) Sleep() {
	c.Calls = append(c.Calls, sleep) 
}

func (c *CountdownOperationSpy) Write(p []byte) (n int, err error) {
	c.Calls = append(c.Calls, write)
	return
}


func TestCountdown(t *testing.T) {
	t.Run("Test countdown output", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		sleeper := &CountdownOperationSpy{}
		Countdown(buffer, sleeper)

		got := buffer.String()
		want := `3
2
1
Go!`
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("Test sleep before print", func(t *testing.T) {
		spy := &CountdownOperationSpy{}

		Countdown(spy, spy)
		want := []string{
			"sleep",
			"write",
			"sleep",
			"write",
			"sleep",
			"write",
			"sleep",
			"write",
		}

		if !reflect.DeepEqual(want, spy.Calls) {
			t.Errorf("want %v got %v", want, spy.Calls)
		}

	})
}

func TestConfigurableSleeper(t *testing.T) {
	sleepTime := 5 * time.Second
	spyTime := &SpyTime{}
	sleeper := ConfigurableSleeper{sleepTime, spyTime.Sleep}
	sleeper.Sleep()

	if spyTime.duration != sleepTime {
		t.Errorf("want sleep for %v, but slept for %v", sleepTime, spyTime.duration)
	}
}


