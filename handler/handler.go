package handler

import (
	"fmt"
	"log"
	"net/http"

	"stash.corp.synacor.com/hack/werewolf/game"
	"stash.corp.synacor.com/hack/werewolf/player"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Disable CORS for testing
	},
}

type Handler struct {
	game     *game.Game
	joinChan chan *player.Player
}

func New() *Handler {
	joinChan := make(chan *player.Player)
	return &Handler{
		game:     game.New(joinChan),
		joinChan: joinChan,
	}
}

func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("BAD upgrade:", err)
		return
	}

	p := player.New(c, h.joinChan)
	go p.Play()
}

// Awoo's only function is to amuse its authors.
func (h *Handler) Awoo(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v awoos\n", r.RemoteAddr)
	fmt.Fprintln(w, "awoooooooooooo")
}
