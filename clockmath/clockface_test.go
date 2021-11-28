package clockface

import (
	"testing"
    "time"
    "math"
)

func simpleTime(hour int, minute int, second int) time.Time {
	return time.Date(1337, time.January, 12, hour, minute, second, 0, time.UTC)
}

func TestSecondHandAtMidnight(t *testing.T) {
	tm := time.Date(1337, time.January, 1, 0, 0, 0, 0, time.UTC)
	want := Point{X: 150, Y: 150-90}
	got := SecondHand(tm)

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSecondHandAt30Seconds(t *testing.T) {
	tm := time.Date(1337, time.January, 1, 0, 0, 30, 0, time.UTC)
	want := Point{X: 150, Y: 150+90}
	got := SecondHand(tm)

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}


func TestSecondsInRadians(t *testing.T) {
	thirtySeconds := time.Date(1337, time.January, 1, 0, 0, 30, 0, time.UTC)
	want := math.Pi
	got := secondsInRadians(thirtySeconds)

	if want != got {
		t.Fatalf("want %v got %v", want, got)
	}
}

func TestSecondHandAtVector(t *testing.T) {
	cases := []struct {
		time time.Time
		point Point
	} {
		{
			simpleTime(0, 0, 30), 
			Point{0, -1},
		},
		{
			simpleTime(0, 0, 45),
			Point{-1, 0},
		},
	}

	for _, c := range cases {
		t.Run(c.time.Format("12:14:04"), func(t *testing.T) {
			got := secondHandPoint(c.time)
			if !roughlyEqualPoint(got, c.point) {
				t.Fatalf("want %v, got %v", c.point, got)
			}
		})
	}
}


func roughlyEqualFloat64(a, b float64) bool {
	const threshold = 1e-7
	return math.Abs(a - b) < threshold
}

func roughlyEqualPoint(p1, p2 Point) bool {
	return roughlyEqualFloat64(p1.X, p2.X) && roughlyEqualFloat64(p1.Y, p2.Y)
}





