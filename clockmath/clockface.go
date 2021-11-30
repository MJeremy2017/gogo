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
	p := secondHandPoint(t)
	p = Point{clockCenterX + secondHandLength * p.X, clockCenterY - secondHandLength * p.Y}
	return p
}

func secondsInRadians(t time.Time) float64 {
	return math.Pi / (30 / float64(t.Second()))
}


func secondHandPoint(t time.Time) Point {
	radian := secondsInRadians(t)
	return Point{math.Sin(radian), math.Cos(radian)}
}