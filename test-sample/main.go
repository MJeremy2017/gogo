package main

import (
	"fmt"
	"sync"
	"time"
	"something/innerpkg"
)

var mutex sync.Mutex
var a int
var wg sync.WaitGroup

type Animal interface {
	Yell() string

	Walk() string
}

type Bird struct {
	name string
	sound string
}

func (b *Bird) Yell() string {
	return b.sound
}

func (b *Bird) Walk() string {
	return ""
}


func GetAnimalSound(animal Animal) string {
	return animal.Yell()
}

func main() {
	count := 1000
	wg.Add(count)
	for i := 0; i < count; i++ {
		go Incr()
	}
	wg.Wait()

	fmt.Printf("a: %d\n", a)

	bird := &Bird{"canical", "gee"}
	sound := GetAnimalSound(bird)
	fmt.Printf("%q sound is %q", bird.name, sound)

}


func Incr() {
	mutex.Lock()
	defer mutex.Unlock()
	time.Sleep(time.Duration(1 * time.Millisecond))
	a++
	wg.Done()
}