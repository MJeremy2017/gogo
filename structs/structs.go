package structs

import "math"

type Rectangle struct {
	Width float64
	Height float64
}

type Circle struct {
	Radius float64
}

type Triangle struct {
	Width float64
	Height float64
}

type Shape interface {
	Area() float64
}

// a method associates to a type

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (t Triangle) Area() float64 {
	return 0.5 * t.Width * t.Height
}

func Perimeter(rec Rectangle) float64 {
	return 2 * (rec.Width + rec.Height)
}
