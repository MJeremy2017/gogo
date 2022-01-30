package some

import (
	"io"
)


func NewCancellatbleReader(rdr io.Reader) io.Reader {
	return rdr
}