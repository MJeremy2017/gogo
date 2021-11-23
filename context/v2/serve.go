package mycontext

import (
	"net/http"
	"context"
	"fmt"
)


type Store interface {
	Fetch(ctx context.Context) (string, error)
}


func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// past on the request context without dealing in the server side
		data, err := store.Fetch(r.Context())
		if err != nil {
			return
		}
		fmt.Fprint(w, data)
	}
}