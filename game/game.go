package game

import (
	"fmt"
	"log"
	"math/rand"

	"stash.corp.synacor.com/hack/werewolf/chanmsg"
	"stash.corp.synacor.com/hack/werewolf/message"
	"stash.corp.synacor.com/hack/werewolf/player"
	"stash.corp.synacor.com/hack/werewolf/role/roleset"

	"github.com/gofrs/uuid"
)

const (
	NotRunning = iota
	Setup      = iota
	Running    = iota
	Finished   = iota
)

type Game struct {
	Players    map[uuid.UUID]*player.Player `json:"-"`
	PlayerList []*player.Player             `json:"players"`
	Roleset    *roleset.Roleset             `json:"roleset"`
	State      int                          `json:"game_state"`
	Phase      int                          `json:"phase"`
	gameChan   chan *chanmsg.Activity       `json:"-"`
}

func New(joinChan chan *player.Player) *Game {
	game := &Game{
		Players:  make(map[uuid.UUID]*player.Player),
		State:    NotRunning,
		Phase:    0,
		gameChan: make(chan *chanmsg.Activity),
	}
	go game.HandleJoins(joinChan)
	go game.HandlePlayerMessage()
	return game
}

func (g *Game) UpdatePlayerList() {
	var l []*player.Player
	for _, p := range g.Players {
		l = append(l, p)
	}
	g.PlayerList = l
}

func (g *Game) AddPlayer(p *player.Player) error {
	if g.State == NotRunning {
		g.State = Setup
	}
	if g.State != Setup {
		return fmt.Errorf("cannot add player: game is not in setup phase")
	}

	log.Printf("new player: %s", p.UUID)
	if len(g.Players) == 0 {
		log.Printf("setting leader: %s", p.UUID)
		p.Leader = true

		// TODO
		g.SetRoleset(roleset.Fiver())
	}
	p.SetChan(g.gameChan)

	g.Players[p.UUID] = p
	g.UpdatePlayerList()

	// let everyone know they have a new friend
	g.Broadcast(message.PlayerList, g.PlayerList)

	// send them an ack/please wait
	if err := p.Message(message.PleaseWait, p); err != nil {
		log.Println(err)
	}

	if g.ShouldStart() {
		g.Start()
	}
	return nil
}

func (g *Game) ShouldStart() bool {
	return g.Roleset != nil && len(g.Players) == len(g.Roleset.Roles)
}

func (g *Game) Start() error {
	if g.State != Setup {
		return fmt.Errorf("cannot start game: game is not in setup phase")
	}

	roleOrder := rand.Perm(len(g.Roleset.Roles))
	for k, v := range g.PlayerList {
		g.Players[v.UUID].Role = g.Roleset.Roles[roleOrder[k]]
		//g.Players[v.UUID].Message(
	}

	g.State = Running
	g.Phase = 1
	//g.Broadcast(message.GameState, g.State)
	g.Broadcast(message.Phase, g.Phase)

	return nil
}

func (g *Game) SetRoleset(r *roleset.Roleset) error {
	if g.State != Setup {
		return fmt.Errorf("cannot set roleset: game is not in 'not running' phase")
	}

	g.Roleset = r
	log.Printf("set roleset to %s", g.Roleset.Name)
	g.Broadcast(message.Roleset, g.Roleset)

	if g.ShouldStart() {
		g.Start()
	}
	return nil
}

/*
func (g *Game) Vote(v *vote.Vote) {
	if !g.Day {
	}
}
*/

func (g *Game) Day() bool {
	return g.State%2 == 1
}

func (g *Game) RemovePlayer(id uuid.UUID) {
	delete(g.Players, id)
	g.UpdatePlayerList()

	// let everyone know they lost a friend
	g.Broadcast(message.PlayerList, g.PlayerList)
}

func (g *Game) Broadcast(title string, payload interface{}) {
	for _, p := range g.Players {
		if err := p.Message(title, payload); err != nil {
			log.Printf("%s: error messaging (%s)", p.UUID, err)
		}
	}
}

func (g *Game) HandleJoins(joinChan chan *player.Player) {
	for {
		p := <-joinChan
		log.Printf("%s: joining", p.UUID)
		g.AddPlayer(p)
	}
}

func (g *Game) HandlePlayerMessage() {
	for {
		activity := <-g.gameChan
		switch activity.Type {
		case chanmsg.Quit:
			log.Printf("%s: quitting", *activity.From)
			g.RemovePlayer(*activity.From)
			/*
				case vote := <-g.VoteChan:
					log.Println(vote)
					//g.Vote(vote)
			*/
		case chanmsg.PlayerList:
			log.Printf("%s: requesting player list", *activity.From)
			g.Players[*activity.From].Message(message.PlayerList, g.PlayerList)
		}
	}
}
