package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Sigafoos/awoo/chanmsg"
	"github.com/Sigafoos/awoo/message"
	"github.com/Sigafoos/awoo/player"
	"github.com/Sigafoos/awoo/role"
	"github.com/Sigafoos/awoo/role/roleset"
	"github.com/Sigafoos/awoo/tally"

	"github.com/gofrs/uuid"
)

const (
	NotRunning = iota
	Setup      = iota
	Running    = iota
	Finished   = iota
)

type Game struct {
	Players          map[uuid.UUID]*player.Player `json:"-"`
	PlayerList       []*player.Player             `json:"players"`
	Roleset          *roleset.Roleset             `json:"roleset"`
	votes            map[*player.Player]uuid.UUID `json:"-"`
	Tally            []*tally.TallyItem           `json:"tally"`
	State            int                          `json:"game_state"`
	Phase            int                          `json:"phase"`
	NightActionQueue []*FingerPoint               `json:"-"`
	gameChan         chan *chanmsg.Activity       `json:"-"`
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

	log.Printf("new player: %s", p.Identifier())
	if len(g.Players) == 0 {
		log.Printf("setting leader: %s", p.Identifier())
		p.Leader = true

		// TODO
		//g.SetRoleset(roleset.Fiver())
		//g.SetRoleset(roleset.Debug())
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

	log.Println("== role selection ==")
	rand.Seed(time.Now().UnixNano())
	roleOrder := rand.Perm(len(g.Roleset.Roles))
	for k, v := range g.PlayerList {
		r := g.Roleset.Roles[roleOrder[k]]
		log.Printf("%s: %s", v.Identifier(), r.Name)
		g.Players[v.UUID].Role = r
		g.Players[v.UUID].Message(message.Role, r)
	}

	log.Printf("== starting game ==")

	g.State = Running
	g.NextPhase()
	g.ProcessStartActionQueue()
	//g.Broadcast(message.GameState, g.State)

	return nil
}

func (g *Game) NextPhase() {
	g.UpdatePlayerList()
	maxes := g.AliveMaxEvils()
	g.NightActionQueue = []*FingerPoint{}
	g.votes = make(map[*player.Player]uuid.UUID)
	g.RebuildTally()

	log.Printf("alive maxes: %v", maxes)
	log.Printf("parity: %v", g.Parity())

	if len(maxes) == 0 {
		g.End(role.Good)
		return
	} else if g.Parity() <= 0 {
		g.End(role.Evil)
		return
	} else if len(g.PlayerList) == 2 {
		g.End(role.Good)
		return
	}

	g.Phase++

	log.Printf("== game is now on phase %v ==", g.Phase)
	g.Broadcast(message.Phase, g.Phase)
}

func (g *Game) Parity() int {
	parity := 0
	for _, p := range g.Players {
		if p.Role.Alive {
			parity += p.Role.Parity
		}
	}
	return parity
}

// RebuildTally calculates the current tally based on individual votes.
// If there are no votes, it essentially clears it.
func (g *Game) RebuildTally() {
	t := make(map[*player.Player][]*player.Player)

	for _, p := range g.PlayerList {
		t[p] = []*player.Player{}
	}

	for from, to := range g.votes {
		t[g.Players[to]] = append(t[g.Players[to]], from)
	}

	list := []*tally.TallyItem{}
	for c, v := range t {
		item := tally.Item(c, v)
		list = append(list, item)
	}
	g.Tally = list

	if g.Day() {
		g.Broadcast(message.Tally, tally.Short(g.Tally))
	}
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

	log.Printf("killed: %s (%s)", revealed.Name, revealed.RoleName)
	g.NextPhase()
	g.Broadcast(message.Targeted, revealed)
	ousted.Message(message.Dead, "")
}

func (g *Game) AliveMaxEvils() []string {
	var maxes []string
	for _, p := range g.PlayerList {
		if p.Role.Parity < 0 {
			maxes = append(maxes, p.Identifier())
		}
	}
	return maxes
}

func (g *Game) End(victor int) {
	log.Printf("== game over. victor: %v ==", victor)
	g.Broadcast(message.Victory, victor)
	g.State = Finished
}

func (g *Game) VotedOut() *player.Player {
	// if an even number, first to half. otherwise, 50%+1
	threshold := len(g.PlayerList) / 2
	if len(g.PlayerList)%2 == 1 {
		threshold++
	}

	for _, item := range g.Tally {
		if len(item.Votes) >= threshold {
			return item.Candidate
		}
	}
	return nil
}

func (g *Game) QueueNightAction(fp *FingerPoint) {
	g.NightActionQueue = append(g.NightActionQueue, fp)
	log.Printf("%v / %v actions", len(g.NightActionQueue), len(g.PlayerList))
	if len(g.NightActionQueue) >= len(g.PlayerList) {
		g.ProcessNightActionQueue()
	}
}

func (g *Game) ProcessStartActionQueue() {
	for _, p := range g.Players {
		result := g.StartAction(p)
		if result.PlayerMessage != "" {
			p.Message(message.Private, result.PlayerMessage)
		}
	}
}

func (g *Game) ProcessNightActionQueue() {
	var deaths []*player.Revealed
	for _, action := range g.NightActionQueue {
		result := g.NightAction(action)
		if result.PlayerMessage != "" {
			action.From.Message(message.Private, result.PlayerMessage)
		}
		if result.Killed != nil {
			deaths = append(deaths, action.To.Reveal())
			action.To.Message(message.Dead, "")
		}
	}

	g.NextPhase()

	if len(deaths) > 0 {
		// TODO maybe more than one death?
		g.Broadcast(message.Targeted, deaths[0])
	}
}

func (g *Game) Day() bool {
	return g.Phase%2 == 1
}

func (g *Game) RemovePlayer(id uuid.UUID) {
	delete(g.Players, id)
	g.UpdatePlayerList()

	// let everyone know they lost a friend
	g.Broadcast(message.PlayerList, g.PlayerList)

	if len(g.PlayerList) == 0 {
		g.State = NotRunning
		g.Phase = 0
	}
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
		log.Printf("%s: joining", p.Identifier())
		if err := g.AddPlayer(p); err != nil {
			p.Message(message.CanNotJoin, "")
		}
	}
}

func (g *Game) HandlePlayerMessage() {
	for {
		activity := <-g.gameChan
		from := g.Players[activity.From]
		to := g.Players[activity.To]

		switch activity.Type {
		case chanmsg.Quit:
			log.Printf("%s: quitting", from.Identifier())
			g.RemovePlayer(activity.From)
		case chanmsg.PlayerList:
			log.Printf("%s: requesting player list", from.Identifier())
			from.Message(message.PlayerList, g.PlayerList)

		case chanmsg.Tally:
			log.Printf("%s: requesting tally", from.Identifier())
			from.Message(message.Tally, tally.Short(g.Tally))

		case chanmsg.SetRoleset:
			log.Printf("%s: setting roleset %s", from.Identifier(), activity.Roleset)
			sets := roleset.List()
			g.SetRoleset(sets[activity.Roleset])

		case chanmsg.Vote:
			log.Printf("%s: voting for %s", from.Identifier(), g.Players[activity.To].Identifier())
			if err := g.Vote(g.Players[activity.From], activity.To); err != nil {

			}

		case chanmsg.NightAction:
			log.Printf("%s: night action submitted", from.Identifier())
			fp := &FingerPoint{
				From: from,
				To:   to,
			}
			g.QueueNightAction(fp)
		}

	}
}
