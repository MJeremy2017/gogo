package some

import (
	"io"
	"context"
)


type readerCtx struct {
	ctx  		context.Context
	delegate	io.Reader
}



func NewCancellableReader(ctx context.Context, rdr io.Reader) io.Reader {
	return rdr
}