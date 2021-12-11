package main

import (
	"log"
	"net/http"
)

type InMemoryPlayerStore struct {

}

func (i *InMemoryPlayerStore) GetPlayerScore(player string) int {
	return 123
}

func main() {
	store := &InMemoryPlayerStore{}
	server := &PlayerServer{store}
	// handler := http.HandlerFunc(PlayerServer)  // cast into type HandlerFunc which has implemented serveHttp method already
	log.Fatal(http.ListenAndServe(":5000", server))
}