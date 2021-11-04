package main

import "testing"

func TestSearch(t *testing.T) {
	t.Run("known word", func(t *testing.T) {
		dict := Dictionary{"test": "this is a test"}

		got, _ := dict.Search("test")
		want := "this is a test"
		assertStrings(t, got, want)
	})

	t.Run("unknown word", func(t *testing.T) {
		dict := Dictionary{"test": "is a test"}

		_, err := dict.Search("something")
		want := ErrKeyNotFound
		if err == nil {
			t.Fatal("expect an error but not find")
		}
		assertError(t, err, want)
	})

}

func TestAdd(t *testing.T) {
	dict := Dictionary{}
	dict.Add("test", "this is a test")
	got, err := dict.Search("test")
	want := "this is a test"
	if err != nil {
		t.Fatal("Does not expect an error! error:\n", err)
	}

	assertStrings(t, got, want)
}

func assertStrings(t testing.TB, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertError(t testing.TB, got error, want error) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}