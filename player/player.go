package player

import (
	"log"

	"github.com/awoo-detat/awoo/chanmsg"
	"github.com/awoo-detat/awoo/message"
	"github.com/awoo-detat/awoo/role"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

type GamePlayer struct {
	UUID     uuid.UUID `json:"uuid"`
	Leader   bool      `json:"leader"`
	Name     string    `json:"name,omitempty"`
	role     *role.Role
	socket   Communicator
	joinChan chan Player
	gameChan chan *chanmsg.Activity
}

type Revealed struct {
	UUID     uuid.UUID `json:"uuid"`
	Name     string    `json:"name,omitempty"`
	RoleName string    `json:"role"`
	Alive    bool      `json:"alive"`
}

func New(socket Communicator, joinChan chan Player) *GamePlayer {
	id, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	p := &GamePlayer{
		UUID:     id,
		socket:   socket,
		role:     &role.Role{},
		joinChan: joinChan,
	}

	if err := p.Message(message.Awoo, p.ID()); err != nil {
		log.Print(err)
	}
	return p
}

func (p *GamePlayer) ID() string {
	return p.UUID.String()
}

func (p *GamePlayer) Identifier() string {
	if p.Name != "" {
		return p.Name
	}
	return p.ID()
}

func (p *GamePlayer) Reveal() *Revealed {
	r := &Revealed{
		Name: p.Name,
		UUID: p.UUID,
	}
	if p.Role().Alive {
		r.Alive = true
	} else {
		r.RoleName = p.Role().Name
	}
	return r
}

// Message wraps the error checking around encoding a message.Message,
// and sends it on the websocket.
func (p *GamePlayer) Message(title string, payload interface{}) error {
	m, err := message.New(title, payload)
	if err != nil {
		return err
	}
	return p.socket.WriteMessage(websocket.TextMessage, m)
}

func (p *GamePlayer) SetChan(c chan *chanmsg.Activity) {
	p.gameChan = c
}

func (p *GamePlayer) Vote(to string) {
	vote := chanmsg.New(chanmsg.Vote, p.ID())
	vote.To = to
	p.gameChan <- vote
}

func (p *GamePlayer) NightAction(to string) {
	action := chanmsg.New(chanmsg.NightAction, p.ID())
	action.To = to
	p.gameChan <- action
}

// Play is the loop that runs for a websocket to communicate between the
// client and server. If websockets are not being used, this will not trigger.
func (p *GamePlayer) Play() {
	defer p.socket.Close()
	for {
		messageType, content, err := p.socket.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("%s: closing connection", p.Identifier())
				break
			}
			log.Printf("%s: websocket error (%s), closing", p.Identifier(), err)
			break
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
			p.gameChan <- chanmsg.New(chanmsg.PlayerList, p.ID())
		} else if m.PollTally {
			p.gameChan <- chanmsg.New(chanmsg.Tally, p.ID())
		} else if m.Roleset != "" {
			activity := chanmsg.New(chanmsg.SetRoleset, p.ID())
			activity.Roleset = m.Roleset
			p.gameChan <- activity
		} else if m.Vote != "" {
			if m.Time == message.TimeDay {
				p.Vote(m.Vote)
			} else if m.Time == message.TimeNight {
				p.NightAction(m.Vote)
			} else {
				p.Message(message.Error, "jcantwell what did you do it is neither day nor night; are we trapped in this eternal twilight together now?")
			}
		} else {
			log.Printf("unknown request from %s: %s", p.Identifier(), content)
		}
	}
}

func (p *GamePlayer) Reconnect(c Communicator) {
	log.Printf("%s: reconnecting", p.Identifier())
	p.socket = c
	p.Message(message.PleaseWait, p)
	go p.Play()
}

func (p *GamePlayer) Role() *role.Role {
	return p.role
}

func (p *GamePlayer) SetLeader() {
	p.Leader = true
}

func (p *GamePlayer) SetName(name string) {
	p.Name = name
}

func (p *GamePlayer) SetRole(r *role.Role) {
	p.role = r
}

func (p *GamePlayer) InGame() bool {
	return p.gameChan != nil
}

func (p *GamePlayer) LeaveGame() {
	log.Printf("%s: leaving game", p.Identifier())
	p.gameChan = nil
}

func (p *GamePlayer) Quit() {
	if err := p.socket.Close(); err != nil {
		log.Printf("error closing channel: %s", err)
	}
	p.gameChan <- chanmsg.New(chanmsg.Quit, p.ID())
	p.LeaveGame()
}
