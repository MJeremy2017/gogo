package main

import (
	"io"
	// "fmt"
	"encoding/json"
)


type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
	league League
}


func (f *FileSystemPlayerStore) GetLeague() League {
	// offset and whence
	return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	player := f.league.Find(name)

	if player != nil {
		return player.Wins
	}
	return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()
	player := league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}
	f.database.Seek(0, 0)
	json.NewEncoder(f.database).Encode(f.league)
}

func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
	database.Seek(0, 0)
	league, _ := NewLeague(database)
	return &FileSystemPlayerStore{
		database: database,
		league: league,
	}
}



