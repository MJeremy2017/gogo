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

func MinuteHand(t time.Time) Point {
	p := minuteHandPoint(t)
	p = Point{clockCenterX + minuteHandLength * p.X, clockCenterY - minuteHandLength * p.Y}
	return p
}

func makeHand(p Point, length float64) Point {
	p = Point{clockCenterX + length * p.X, clockCenterY - length * p.Y}
	return p
}

func secondsInRadians(t time.Time) float64 {
	return math.Pi / (30 / float64(t.Second()))
}

func minutesInRadians(t time.Time) float64 {
	minute := float64(t.Minute())
	second := float64(t.Second())
	return minute * math.Pi / 30.0 + second * math.Pi / (60 * 30.0)
}

func hoursInRadians(t time.Time) float64 {
	hour := float64(t.Hour() % 12) 
	return hour * math.Pi/6 + minutesInRadians(t)/12
}

func secondHandPoint(t time.Time) Point {
	radian := secondsInRadians(t)
	return angleToPoint(radian)
}

func minuteHandPoint(t time.Time) Point {
	radian := minutesInRadians(t)
	return angleToPoint(radian)
}

func angleToPoint(angle float64) Point {
	return Point{math.Sin(angle), math.Cos(angle)}
}