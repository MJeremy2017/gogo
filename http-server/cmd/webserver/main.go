package main

import (
	"fmt"
	"time"
)

const dbFileName = "game.db.json"

type A interface {
	GetName() string
	GetHeight() float32
}

type AB interface {
	GetName() string
	GetAge() int32
}

type Ali struct {
}

func (a Ali) GetName() string {
	return "ali"
}

func (a Ali) GetAge() int32 {
	return 12
}

func (a Ali) GetHeight() float32 {
	return 12.3
}

func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}


func main() {
	fmt.Printf("time %s", time.Now().UTC())
	//store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	//dummyGame := &poker.GameSpy{}
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer close()
	//
	//server, _ := poker.NewPlayerServer(store, dummyGame)
	//// handler := http.HandlerFunc(PlayerServer)  // cast into type HandlerFunc which has implemented serveHttp method already
	//log.Println("listen on port 5000")
	//log.Fatal(http.ListenAndServe(":5000", server))
}
