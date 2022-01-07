package main

import (
	"fmt"
	"os"
	"log"
	"server/poker"
	"server/poker_test"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")
	dummyAlerter := &poker_test.SpyBlindAlerter{}
	poker.NewCLI(store, os.Stdin, dummyAlerter).PlayPoker()
}