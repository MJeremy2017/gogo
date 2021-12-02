package clockface

import (
	"testing"
	"encoding/xml"
    "time"
    "math"
    "bytes"
)

type SVG struct {
	XMLName xml.Name `xml:"svg"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	ViewBox string   `xml:"viewBox,attr"`
	Version string   `xml:"version,attr"`
	Circle  Circle   `xml:"circle"`
	Line    []Line   `xml:"line"`
}

type Line struct {
	X1 float64 `xml:"x1,attr"`
	Y1 float64 `xml:"y1,attr"`
	X2 float64 `xml:"x2,attr"`
	Y2 float64 `xml:"y2,attr"`
}

type Circle struct {
	Cx float64 `xml:"cx,attr"`
	Cy float64 `xml:"cy,attr"`
	R  float64 `xml:"r,attr"`
}


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

func TestMinutesInRadians(t *testing.T) {
	cases := []struct {
		time time.Time
		angle float64
	} {
		{simpleTime(0, 30, 0), math.Pi},
		{simpleTime(0, 0, 7), 7 * math.Pi / (30*60)},
	}
	
	for _, tc := range cases {
		t.Run(tc.time.Format("12:14:04"), func(t *testing.T) {
			got := minutesInRadians(tc.time)
			want := tc.angle
			if !roughlyEqualFloat64(want, got) {
				t.Fatalf("want %v got %v", want, got)
			}
		})

	}
}

func TestHoursInRadians(t *testing.T) {
	cases := []struct {
		time time.Time
		angle float64
	} {
		{simpleTime(6, 0, 0), math.Pi},
		{simpleTime(0, 1, 30), 90 * math.Pi / (6*60*60)},
	}
	
	for _, tc := range cases {
		t.Run(tc.time.Format("12:14:04"), func(t *testing.T) {
			got := hoursInRadians(tc.time)
			want := tc.angle
			if !roughlyEqualFloat64(want, got) {
				t.Fatalf("want %v got %v", want, got)
			}
		})

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

func TestMinuteHandAtVector(t *testing.T) {
	cases := []struct {
		time time.Time
		point Point
	} {
		{
			simpleTime(0, 30, 0), 
			Point{0, -1},
		},
		{
			simpleTime(0, 45, 0),
			Point{-1, 0},
		},
	}

	for _, c := range cases {
		t.Run(c.time.Format("12:14:04"), func(t *testing.T) {
			got := minuteHandPoint(c.time)
			if !roughlyEqualPoint(got, c.point) {
				t.Fatalf("want %v, got %v", c.point, got)
			}
		})
	}
}

func TestHourHandAtVector(t *testing.T) {
	cases := []struct {
		time time.Time
		point Point
	} {
		{
			simpleTime(6, 0, 0), 
			Point{0, -1},
		},
		{
			simpleTime(21, 0, 0),
			Point{-1, 0},
		},
	}

	for _, c := range cases {
		t.Run(c.time.Format("12:14:04"), func(t *testing.T) {
			got := hourHandPoint(c.time)
			if !roughlyEqualPoint(got, c.point) {
				t.Fatalf("want %v, got %v", c.point, got)
			}
		})
	}
}

func TestSVGWriterSecondHand(t *testing.T) {
	cases := []struct {
		time time.Time
		line Line
	}{
		{
			simpleTime(0, 0, 0),
			Line{150, 150, 150, 60},
		},
		{
			simpleTime(0, 0, 30),
			Line{150, 150, 150, 240},
		},
	}

	for _, c := range cases {
		t.Run(c.time.Format("12:12:12"), func(t *testing.T) {
			b := bytes.Buffer{}
			SVGWriter(&b, c.time)

			svg := SVG{}
			xml.Unmarshal(b.Bytes(), &svg)

			if !containsLine(c.line, svg.Line) {
				t.Errorf("Expected to find the second hand line %+v, in the svg line %+v", c.line, svg.Line)
			}
		})
	}
}

func TestSVGWriterHourHand(t *testing.T) {
	cases := []struct {
		time time.Time
		line Line
	}{
		{
			simpleTime(6, 0, 0),
			Line{150, 150, 150, 200},
		},
	}

	for _, c := range cases {
		t.Run(c.time.Format("12:12:12"), func(t *testing.T) {
			b := bytes.Buffer{}
			SVGWriter(&b, c.time)

			svg := SVG{}
			xml.Unmarshal(b.Bytes(), &svg)

			if !containsLine(c.line, svg.Line) {
				t.Errorf("Expected to find the hour hand line %+v, in the svg line %+v", c.line, svg.Line)
			}
		})
	}
}

func containsLine(l Line, ls []Line) bool {
	for _, line := range ls {
		if line == l {
			return true
		}
	}
	return false
}


func roughlyEqualFloat64(a, b float64) bool {
	const threshold = 1e-7
	return math.Abs(a - b) < threshold
}

func roughlyEqualPoint(p1, p2 Point) bool {
	return roughlyEqualFloat64(p1.X, p2.X) && roughlyEqualFloat64(p1.Y, p2.Y)
}





