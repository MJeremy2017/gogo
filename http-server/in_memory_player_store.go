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

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		make(map[string]int),
		sync.Mutex{},
	}
}
