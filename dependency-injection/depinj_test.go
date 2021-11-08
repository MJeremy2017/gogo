package depinj

import (
	"testing"
	"bytes"
)

func TestGreet(t *testing.T) {
	name := "Jeremy"
	buffer := bytes.Buffer{}
	Greet(&buffer, name)
	
	got := buffer.String()
	want := "Hello, Jeremy"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}