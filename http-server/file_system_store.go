package main

import (
	"io"
	"fmt"
)


type FileSystemPlayerStore struct {
	database io.ReadSeeker
}


func (f *FileSystemPlayerStore) GetLeague() []Player {
	// offset and whence
	f.database.Seek(0, 0)
	league, err := NewLeague(f.database)
	if err != nil {
		fmt.Printf("%v", err)
		return nil
	}
	return league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	return 22
}