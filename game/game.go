package game

import (
	"fmt"
	"log"
	"math/rand"

	"stash.corp.synacor.com/hack/werewolf/chanmsg"
	"stash.corp.synacor.com/hack/werewolf/message"
	"stash.corp.synacor.com/hack/werewolf/player"
	"stash.corp.synacor.com/hack/werewolf/role"
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
	Players    map[uuid.UUID]*player.Player        `json:"-"`
	PlayerList []*player.Player                    `json:"players"`
	Roleset    *roleset.Roleset                    `json:"roleset"`
	votes      map[*player.Player]uuid.UUID        `json:"-"`
	Tally      map[*player.Player][]*player.Player `json:"tally"`
	State      int                                 `json:"game_state"`
	Phase      int                                 `json:"phase"`
	gameChan   chan *chanmsg.Activity              `json:"-"`
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
		if p.Role.Name == "" || p.Role.Alive {
			l = append(l, p)
		}
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
		r := g.Roleset.Roles[roleOrder[k]]
		g.Players[v.UUID].Role = r
		g.Players[v.UUID].Message(message.Role, r)
	}

	log.Printf("== starting game ==")

	g.State = Running
	g.NextPhase()
	//g.Broadcast(message.GameState, g.State)

	return nil
}

func (g *Game) NextPhase() {
	g.votes = make(map[*player.Player]uuid.UUID)
	g.Phase++
	g.RebuildTally()
	g.Broadcast(message.Phase, g.Phase)
}

// RebuildTally calculates the current tally based on individual votes.
// If there are no votes, it essentially clears it.
func (g *Game) RebuildTally() {
	tally := make(map[*player.Player][]*player.Player)
	for _, p := range g.PlayerList {
		tally[p] = []*player.Player{}
	}

	for from, to := range g.votes {
		tally[g.Players[to]] = append(tally[g.Players[to]], from)
	}

	g.Tally = tally
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

func (g *Game) Vote(from *player.Player, to uuid.UUID) error {
	if !g.Day() {
		err := fmt.Errorf("vote failed; not day")
		log.Println(err)
		return err
	}

	log.Printf("from %s", from.Name)
	g.votes[from] = to
	g.RebuildTally()

	ousted := g.VotedOut()
	if ousted != nil {
		g.EndDay(ousted)
	}
	return nil
}

func (g *Game) EndDay(ousted *player.Player) {
	// if the player died, regenerate our list
	if !ousted.Role.Kill() {
		g.UpdatePlayerList()
	}
	revealed := ousted.Reveal()

	log.Printf("%+v", revealed)
	g.Broadcast(message.Targeted, revealed)

	maxes := g.AliveMaxEvils()
	if maxes == 0 { // TODO HERE HERE HERE vvv
		g.End(role.Good)
		//} else if g.Parity() <= 0 {
		//g.End(role.Evil)
	} // TODO hunter victory
	// here I'd end the day
}

func (g *Game) AliveMaxEvils() int {
	maxes := 0
	for _, p := range g.PlayerList {
		if p.Role.Parity < 0 {
			maxes++
		}
	}
	return maxes
}

func (g *Game) End(victor int) {
	g.Broadcast(message.Victory, victor)
	g.State = Finished
}

func (g *Game) VotedOut() *player.Player {
	// if an even number, first to half. otherwise, 50%+1
	threshold := len(g.PlayerList) / 2
	if len(g.PlayerList)%2 == 1 {
		threshold++
	}

	for p, votes := range g.Tally {
		if len(votes) >= threshold {
			return p
		}
	}
	return nil
}

func (g *Game) Day() bool {
	return g.Phase%2 == 1
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
			log.Printf("%s: quitting", activity.From)
			g.RemovePlayer(activity.From)
			/*
				case vote := <-g.VoteChan:
					log.Println(vote)
					//g.Vote(vote)
			*/
		case chanmsg.PlayerList:
			log.Printf("%s: requesting player list", activity.From)
			g.Players[activity.From].Message(message.PlayerList, g.PlayerList)
		case chanmsg.Vote:
			log.Printf("%s: voting for %s", activity.From, activity.To)
			g.Vote(g.Players[activity.From], activity.To)
		}
	}
}
