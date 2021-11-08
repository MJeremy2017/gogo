package depinj

import (
	"fmt"
	"io"
)


func Greet(writer io.Writer, name string) {
	// Fprintf writes to the writer instead of printing to stdout
	fmt.Fprintf(writer, "Hello, %s", name)
}