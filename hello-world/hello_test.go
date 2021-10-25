package main

import "testing"

func TestHello(t *testing.T) {
	assertCode := func(t testing.TB, got string, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}

	t.Run("say hello to people", func(t *testing.T) {
		got := Hello("Jacob")
		want := "Hello, Jacob"

		assertCode(t, got, want)
	})

	t.Run("hello with empty input", func(t *testing.T) {
		got := Hello("")
		want := "Hello, World"

		assertCode(t, got, want)
	})
}