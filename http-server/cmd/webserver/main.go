package main

import (
	"log"
	"net/http"
	"server/poker"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	dummyGame := &poker.GameSpy{}
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	server, _ := poker.NewPlayerServer(store, dummyGame)
	// handler := http.HandlerFunc(PlayerServer)  // cast into type HandlerFunc which has implemented serveHttp method already
	log.Println("listen on port 5000")
	log.Fatal(http.ListenAndServe(":5000", server))
}
