package poker

import (
	"os"
)


type tape struct {
	file *os.File
}


func (t tape) Write(data []byte) (n int, err error) {
	t.file.Truncate(0)
	t.file.Seek(0, 0)
	return t.file.Write(data)
}