package main

import "fmt"

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
}