package main

import "fmt"

type ST struct {
	name string
	arr []string
}

func modify(val ST) {
	val.name = "updated"
	val.arr[0] = "updated"
}

func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}
	
	prefix := greetingPrefix(language)
	return prefix + name
}

func greetingPrefix(language string) string {
	var prefix string
	const spanishPrefix = "Hola, "
	const englishPrefix = "Hello, "
	const frenchPrefix = "Bonjour, "

	switch language {
		case "Spanish":
			prefix = spanishPrefix
		case "French":
			prefix = frenchPrefix
		default:
			prefix = englishPrefix
	}
	return prefix
}

func main() {
	fmt.Println(Hello("World", ""))

	st := ST{
		name: "Alfred",
		arr: []string{"Alfred"},
	}

	fmt.Printf("before %v\n", st)
	// only the copy of address of the arr is passed in
	modify(st)
	fmt.Printf("after %v", st)
}