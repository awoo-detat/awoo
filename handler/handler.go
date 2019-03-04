package handler

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/awoo-detat/awoo/game"
	"github.com/awoo-detat/awoo/player"

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
	ips      map[string]*player.Player
}

func New() *Handler {
	joinChan := make(chan *player.Player)
	return &Handler{
		game:     game.New(joinChan),
		joinChan: joinChan,
		ips:      make(map[string]*player.Player),
	}
}

func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("BAD upgrade:", err)
		return
	}

	// handle reconnects
	ip := h.ipFromRequest(r)
	if p, exists := h.ips[ip]; exists && p.Name != "" {
		p.Reconnect(c)
	} else {
		p := player.New(c, h.joinChan)
		query := r.URL.Query()
		if _, set := query["noreconnect"]; !set {
			h.ips[ip] = p
		} else {
			// If the IP is in the map it could still be "new", as in the player
			// doesn't have a name set and never joined the game. In that case,
			// treat them like a completely new user and don't save the IP.
			log.Printf("%s: not saving IP", p.Identifier())
			delete(h.ips, ip)
		}
		go p.Play()
	}
}

// Awoo's only function is to amuse its authors.
func (h *Handler) Awoo(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v awoos\n", r.RemoteAddr)
	fmt.Fprintln(w, "awoooooooooooo")
}

func (h *Handler) ipFromRequest(r *http.Request) string {
	// This bit hasn't been tested. It also doesn't handle case differences.
	ips := r.Header.Get("X-Forwarded-For")
	if ips != "" {
		list := strings.Split(ips, ", ")
		return list[0]
	}

	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return ""
}
