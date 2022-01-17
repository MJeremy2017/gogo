package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"html/template"
)

const jsonContentType = "application/json"

type PlayerStore interface {

	GetPlayerScore(player string) int

	RecordWin(name string)

	GetLeague() League
}

// server takes in a player store interface
type PlayerServer struct {
	store        PlayerStore
	http.Handler // embedding, now server can have Handler's methods
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)
	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playerHandler))
	router.Handle("/game", http.HandlerFunc(p.game))

	p.Handler = router
	return p
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.store.GetLeague())
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) playerHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) game(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("cmd/webserver/game.html")

	if err != nil {
		http.Error(w, fmt.Sprintf("problem loading template %s", err.Error()),
			http.StatusInternalServerError) 
		return
	}

	tmpl.Execute(w, nil)

	w.WriteHeader(http.StatusOK)

}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
	return
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprintf(w, strconv.Itoa(score))
	return

}
