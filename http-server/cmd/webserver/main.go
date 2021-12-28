package main

import (
	"log"
	"net/http"
	"os"
	"poker"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem openning file %s %v", dbFileName, err)
	}
	store, err := poker.NewFileSystemPlayerStore(db)
	if err != nil {
		log.Fatalf("got error %v", err)
	}
	server := poker.NewPlayerServer(store)
	// handler := http.HandlerFunc(PlayerServer)  // cast into type HandlerFunc which has implemented serveHttp method already
	log.Println("listen on port 5000")
	log.Fatal(http.ListenAndServe(":5000", server))
}
