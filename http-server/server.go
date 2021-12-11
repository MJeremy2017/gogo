package main

import (
	"fmt"
	"net/http"
	"strings"
	"strconv"
)

type PlayerStore interface {
	GetPlayerScore(player string) int
}

type PlayerServer struct {
	store PlayerStore
}


func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	score := p.store.GetPlayerScore(player)
	fmt.Fprintf(w, strconv.Itoa(score))
	return
}