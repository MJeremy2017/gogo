package main

import (
	"testing"
	"io/ioutil"
	"os"
	"log"
)


func TestFileSystemStore(t *testing.T) {
	t.Run("return sorted league", func(t *testing.T) {
		data := `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 22}]`

		database, cleanDatabase := createTempFile(t, data)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)


		got := store.GetLeague()
		want := []Player{
			{"Chris", 22},
			{"Cleo", 10},
		}
		assertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		data := `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 22}]`

		database, cleanDatabase := createTempFile(t, data)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		got := store.GetPlayerScore("Chris")
		want := 22
		assertScoreEqual(t, got, want)
	})

	t.Run("store wins for existing players", func(t *testing.T) {
		data := `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 22}]`

		database, cleanDatabase := createTempFile(t, data)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		store.RecordWin("Chris")

		got := store.GetPlayerScore("Chris")
		want := 23
		assertScoreEqual(t, got, want)
	})

	t.Run("store wins for new players", func(t *testing.T) {
		data := `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 22}]`

		database, cleanDatabase := createTempFile(t, data)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		store.RecordWin("Pepper")

		got := store.GetPlayerScore("Pepper")
		want := 1
		assertScoreEqual(t, got, want)
	})

	t.Run("no error with empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)
	})
}

func createTempFile(t testing.TB, data string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}
	tmpfile.Write([]byte(data))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}
	return tmpfile, removeFile
}

func assertScoreEqual(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		log.Fatalf("didn't expect error but got %v", err)
	}
}

