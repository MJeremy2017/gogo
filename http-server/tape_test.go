package main

import (
	"testing"
	"io/ioutil"
)

func TestTape_Writer(t *testing.T) {
	want := "abc"
	file, clean := createTempFile(t, "1234")
	defer clean()

	tape := &tape{file}
	tape.Write([]byte(want))

	file.Seek(0, 0)
	newFileContents, _ := ioutil.ReadAll(file)

	got := string(newFileContents)
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}