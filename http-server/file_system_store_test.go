package main

import (
	"testing"
	"io/ioutil"
	"io"
	"os"
)


func TestFileSystemStore(t *testing.T) {

	t.Run("league from a reader", func(t *testing.T) {
		data := `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 22}]`

		database, cleanDatabase := createTempFile(t, data)
		defer cleanDatabase()

		store := NewFileSystemPlayerStore(database)

		got := store.GetLeague()
		want := []Player{
			{"Cleo", 10},
			{"Chris", 22},
		}
		assertLeague(t, got, want)

		// read a second time
		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		data := `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 22}]`

		database, cleanDatabase := createTempFile(t, data)
		defer cleanDatabase()

		store := NewFileSystemPlayerStore(database)
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

		store := NewFileSystemPlayerStore(database)
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

		store := NewFileSystemPlayerStore(database)
		store.RecordWin("Pepper")

		got := store.GetPlayerScore("Pepper")
		want := 1
		assertScoreEqual(t, got, want)
	})


}

func createTempFile(t testing.TB, data string) (io.ReadWriteSeeker, func()) {
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