package clockface

import (
	"time"
	"math"
)

// A Point represents a two dimensional Cartesian coordinate.
type Point struct {
	X float64
	Y float64
}

// SecondHand is the unit vector of the second hand of an analogue clock at time `t`.
// represented as a Point.
func SecondHand(t time.Time) Point {
	radius := secondsInRadians(t)
	return Point{150 + math.Sin(radius), 150 - 90 * math.Cos(radius)}
}

func secondsInRadians(t time.Time) float64 {
	return math.Pi / (30 / float64(t.Second()))
}


func secondHandPoint(t time.Time) Point {
	radian := secondsInRadians(t)
	return Point{math.Sin(radian), math.Cos(radian)}
}