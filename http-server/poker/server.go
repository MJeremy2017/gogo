package poker

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"log"
)

const jsonContentType = "application/json"
const htmlTemplatePath = "game.html"

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type PlayerStore interface {
	GetPlayerScore(player string) int

	RecordWin(name string)

	GetLeague() League
}

// PlayerServer server takes in a player store interface
type PlayerServer struct {
	store        PlayerStore
	http.Handler // embedding, now server can have Handler's methods
	template     *template.Template
	game         Game
}


type playerServerWS struct {
	*websocket.Conn
}


func (w *playerServerWS) Write(p []byte) (n int, err error) {
	err = w.WriteMessage(websocket.TextMessage, p)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func NewPlayerServerWS(w http.ResponseWriter, r *http.Request) *playerServerWS {
	conn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("problem upgrading connection to WebSockets %v\n", err)
	}

	return &playerServerWS{conn}
}

func (w *playerServerWS) WaitForMsg() string {
	_, msg, err := w.ReadMessage()
	if err != nil {
		log.Printf("error reading from websocket %v\n", err)
	}
	return string(msg)
}


func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {
	p := new(PlayerServer)

	tmpl, err := template.ParseFiles(htmlTemplatePath)

	if err != nil {
		msg := fmt.Errorf("problem loading template path %s %v", htmlTemplatePath, err)
		return nil, msg
	}

	p.store = store
	p.template = tmpl
	p.game = game

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playerHandler))
	router.Handle("/game", http.HandlerFunc(p.playGame))
	router.Handle("/ws", http.HandlerFunc(p.webSocket))

	p.Handler = router
	return p, nil
}

func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := NewPlayerServerWS(w, r)

	numberOfPlayersMsg := ws.WaitForMsg()
	numberOfPlayers, _ := strconv.Atoi(numberOfPlayersMsg)
	p.game.Start(numberOfPlayers, ws)

	winnerMsg := ws.WaitForMsg()
	p.game.Finish(winnerMsg)
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.store.GetLeague())
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

func (p *PlayerServer) playGame(w http.ResponseWriter, r *http.Request) {
	p.template.Execute(w, nil)
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
