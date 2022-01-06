package main

import (
	"log"
	"net/http"
	"server/poker"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	server := poker.NewPlayerServer(store)
	// handler := http.HandlerFunc(PlayerServer)  // cast into type HandlerFunc which has implemented serveHttp method already
	log.Println("listen on port 5000")
	log.Fatal(http.ListenAndServe(":5000", server))
}
