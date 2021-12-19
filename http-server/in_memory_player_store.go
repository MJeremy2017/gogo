package main

import "sync"

type InMemoryPlayerStore struct {
	store map[string]int
	lock sync.Mutex
}

func (i *InMemoryPlayerStore) RecordWin(player string) {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.store[player]++
	return
}

func (i *InMemoryPlayerStore) GetPlayerScore(player string) int {
	return i.store[player]
}

func (i *InMemoryPlayerStore) GetLeague() []Player {
	league := []Player{}
	for key, value := range i.store {
		player := Player{key, value}
		league = append(league, player)
	}
	return league
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		make(map[string]int),
		sync.Mutex{},
	}
}
