package main

import (
	"log"
	"net/http"
)


func main() {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	// handler := http.HandlerFunc(PlayerServer)  // cast into type HandlerFunc which has implemented serveHttp method already
	log.Println("listen on port 5000")
	log.Fatal(http.ListenAndServe(":5000", server))
}
