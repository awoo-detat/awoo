package player

import (
	//"fmt"
	"log"

	"github.com/Sigafoos/awoo/chanmsg"
	"github.com/Sigafoos/awoo/message"
	"github.com/Sigafoos/awoo/role"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	UUID     uuid.UUID  `json:"uuid"`
	Leader   bool       `json:"leader"`
	Name     string     `json:"name,omitempty"`
	Role     *role.Role `json:"-"`
	socket   Communicator
	joinChan chan *Player
	gameChan chan *chanmsg.Activity
}

type Revealed struct {
	UUID     uuid.UUID `json:"uuid"`
	Name     string    `json:"name,omitempty"`
	RoleName string    `json:"role"`
	Alive    bool      `json:"alive"`
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

func (p *Player) Identifier() string {
	if p.Name != "" {
		return p.Name
	}
	return p.UUID.String()
}

func (p *Player) Reveal() *Revealed {
	r := &Revealed{
		Name: p.Name,
		UUID: p.UUID,
	}
	if p.Role.Alive {
		r.Alive = true
	} else {
		r.RoleName = p.Role.Name
	}
	return r
}

// Message wraps the error checking around encoding a message.Message,
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

func (p *Player) Vote(to uuid.UUID) {
	vote := chanmsg.New(chanmsg.Vote, p.UUID)
	vote.To = to
	p.gameChan <- vote
}

func (p *Player) NightAction(to uuid.UUID) {
	action := chanmsg.New(chanmsg.NightAction, p.UUID)
	action.To = to
	p.gameChan <- action
}

// Play is the loop that runs for a websocket to communicate between the
// client and server. If websockets are not being used, this will not trigger.
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
		if m.JoinName != "" {
			p.Name = m.JoinName

			// I don't love that this happens here.
			// TODO make that a separate request
			p.joinChan <- p
		} else if m.PollPlayerList {
			p.gameChan <- chanmsg.New(chanmsg.PlayerList, p.UUID)
		} else if m.PollTally {
			p.gameChan <- chanmsg.New(chanmsg.Tally, p.UUID)
		} else if m.Roleset != "" {
			activity := chanmsg.New(chanmsg.SetRoleset, p.UUID)
			activity.Roleset = m.Roleset
			p.gameChan <- activity
		} else if m.Vote != "" {
			to, err := uuid.FromString(m.Vote)
			if err != nil {
				p.Message(message.Error, err)
				continue
			}
			if m.Time == message.TimeDay {
				p.Vote(to)
			} else if m.Time == message.TimeNight {
				p.NightAction(to)
			} else {
				p.Message(message.Error, "jcantwell what did you do it is neither day nor night; are we trapped in this eternal twilight together now?")
			}
		} else {
			log.Printf("unknown request from %s: %s", p.Identifier(), content)
		}
	}
}

func (p *Player) Quit() {
	if err := p.socket.Close(); err != nil {
		log.Printf("error closing channel: %s", err)
	}
	p.gameChan <- chanmsg.New(chanmsg.Quit, p.UUID)
}
