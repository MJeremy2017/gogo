package main

import "fmt"

func Hello(name string) string {
	if name == "" {
		name = "World"
	}
	const prefix = "Hello, "
	return prefix + name
}

func main() {
	fmt.Println(Hello("World"))
}