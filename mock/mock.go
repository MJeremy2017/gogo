package main

import (
	"io"
	"fmt"
	"os"
)

func Countdown(writer io.Writer) {
	fmt.Fprint(writer, "3")
}

func main() {
	Countdown(os.Stdout)
}