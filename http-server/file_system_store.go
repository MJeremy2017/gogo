package main

import (
	"io"
	"fmt"
	"encoding/json"
)


type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
}


func (f *FileSystemPlayerStore) GetLeague() League {
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
	league := f.GetLeague()
	player := league.Find(name)

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
		league = append(league, Player{name, 1})
	}
	f.database.Seek(0, 0)
	json.NewEncoder(f.database).Encode(league)
}