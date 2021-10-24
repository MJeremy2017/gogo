package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello("Jacob")
	want := "Hello, Jacob"

	if got != want {
		t.Errorf("got %q wang %q", got, want)
	}
}