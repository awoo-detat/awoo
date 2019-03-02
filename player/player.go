package player

import (
	//"fmt"
	"log"

	"stash.corp.synacor.com/hack/werewolf/chanmsg"
	"stash.corp.synacor.com/hack/werewolf/message"
	"stash.corp.synacor.com/hack/werewolf/role"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	UUID     uuid.UUID  `json:"uuid"`
	Leader   bool       `json:"leader"`
	Name     string     `json:"name,omitempty"`
	Role     *role.Role `json:"role"`
	socket   Communicator
	joinChan chan *Player
	gameChan chan *chanmsg.Activity
}

func New(socket Communicator, joinChan chan *Player) *Player {
	id, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	p := &Player{
		UUID:     id,
		socket:   socket,
		Role:     &role.Role{},
		joinChan: joinChan,
	}

	if err := p.Message(message.Awoo, p.UUID.String()); err != nil {
		log.Print(err)
	}
	return p
}

// Message wraps the error checking around encododing a message.Message,
// and sends it on the websocket.
func (p *Player) Message(title string, payload interface{}) error {
	m, err := message.New(title, payload)
	if err != nil {
		return err
	}
	return p.socket.WriteMessage(websocket.TextMessage, m)
}

func (p *Player) SetChan(c chan *chanmsg.Activity) {
	p.gameChan = c
}

func (p *Player) Play() {
	for {
		messageType, content, err := p.socket.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				p.Quit()
				break
			}
			log.Printf("websocket error: %s", err)
		}

		if messageType == websocket.BinaryMessage {
			// this is not ideal.
			log.Panicln("got a binary message...?")
			continue
		}

		m := message.Decode(content)
		log.Printf("%s: %+v", p.UUID, m)
		if m.JoinName != "" {
			p.Name = m.JoinName

			// I don't love that this happens here.
			// TODO make that a separate request
			p.joinChan <- p
		} else if m.PollPlayerList {
			p.gameChan <- chanmsg.New(chanmsg.PlayerList, &p.UUID)
		} else {
			log.Printf("unknown request: %s", content)
		}
	}
}

func (p *Player) Quit() {
	if err := p.socket.Close(); err != nil {
		log.Printf("error closing channel: %s", err)
	}
	p.gameChan <- chanmsg.New(chanmsg.Quit, &p.UUID)
}
