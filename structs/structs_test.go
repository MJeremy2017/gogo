package structs

import (
	"testing"
	"math"
)

func TestPerimeter(t *testing.T) {
	rec := Rectangle{10.0, 10.0}
	got := Perimeter(rec)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}


func TestArea(t *testing.T) {
	areaTests := []struct {
		name string
		shape Shape
		want float64
	} {
		{name: "Rectangle", shape: Rectangle{Width: 12, Height: 3}, want: 36},
		{name: "Circle", shape: Circle{Radius: 10.0}, want: 314.16},
		{name: "Triangle", shape: Triangle{Width: 12, Height: 6}, want: 36},
	}

	for _, test := range areaTests {
		t.Run(test.name, func(t *testing.T) {
			got := math.Round(test.shape.Area()*100)/100
			if got != test.want {
				t.Errorf("%#v got %.2f want %.2f", test.shape, got, test.want)
			}
		})
	}
	
}
