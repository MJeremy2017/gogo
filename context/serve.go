package mycontext

import (
	"net/http"
	"fmt"
)


type Store interface {
	Fetch() string
	Cancel()
}


func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		data := make(chan string, 1)

		go func() {
			data <- store.Fetch()
		}()

		// select among 2 goroutines
		select {
		case d := <-data:
			fmt.Fprint(w, d)
		case <-ctx.Done():
			store.Cancel()
		}

	}
}