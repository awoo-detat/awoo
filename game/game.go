package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/awoo-detat/awoo/chanmsg"
	"github.com/awoo-detat/awoo/message"
	"github.com/awoo-detat/awoo/player"
	"github.com/awoo-detat/awoo/role"
	"github.com/awoo-detat/awoo/role/roleset"
	"github.com/awoo-detat/awoo/tally"
)

const (
	NotRunning = iota
	Setup      = iota
	Running    = iota
	Finished   = iota
)

type Game struct {
	Players          map[string]player.Player `json:"-"`
	PlayerList       []player.Player          `json:"players"`
	Roleset          *roleset.Roleset         `json:"roleset"`
	Tally            []*tally.TallyItem       `json:"tally"`
	State            int                      `json:"game_state"`
	Phase            int                      `json:"phase"`
	NightActionQueue []*FingerPoint           `json:"-"`
	votes            map[player.Player]string
	gameChan         chan *chanmsg.Activity
}

func New(joinChan chan player.Player) *Game {
	game := &Game{
		Players:  make(map[string]player.Player),
		State:    NotRunning,
		Phase:    0,
		gameChan: make(chan *chanmsg.Activity),
	}
	go game.HandleJoins(joinChan)
	go game.HandlePlayerMessage()
	return game
}

func (g *Game) UpdatePlayerList() {
	var l []player.Player
	for _, p := range g.Players {
		if p.Role().Name == "" || p.Role().Alive {
			l = append(l, p)
		}
	}
	g.PlayerList = l
}

func (g *Game) AddPlayer(p player.Player) error {
	if g.State == NotRunning {
		g.State = Setup
	}
	if g.State != Setup {
		return fmt.Errorf("cannot add player: game is not in setup phase")
	}

	log.Printf("new player: %s", p.Identifier())
	if len(g.Players) == 0 {
		log.Printf("setting leader: %s", p.Identifier())
		p.SetLeader()
	}
	p.SetChan(g.gameChan)

	g.Players[p.ID()] = p
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
	if g.Roleset != nil {
		log.Printf("%v/%v players", len(g.Players), len(g.Roleset.Roles))
	}
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
		g.Players[v.ID()].SetRole(r)
		g.Players[v.ID()].Message(message.Role, r)
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
	g.votes = make(map[player.Player]string)
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
		if p.Role().Alive {
			parity += p.Role().Parity
		}
	}
	return parity
}

// RebuildTally calculates the current tally based on individual votes.
// If there are no votes, it essentially clears it.
func (g *Game) RebuildTally() {
	t := make(map[player.Player][]player.Player)

	for _, p := range g.PlayerList {
		t[p] = []player.Player{}
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
	} else if len(r.Roles) < len(g.Players) {
		return fmt.Errorf("roleset %s only has %v roles; %v players in lobby", r.Name, len(r.Roles), len(g.Players))
	}

	g.Roleset = r
	g.Broadcast(message.Roleset, g.Roleset)

	if g.ShouldStart() {
		g.Start()
	}
	return nil
}

func (g *Game) Vote(from player.Player, to string) error {
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

func (g *Game) EndDay(ousted player.Player) {
	// if the player died, regenerate our list
	if !ousted.Role().Kill() {
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
		if p.Role().IsMaxEvil() {
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

func (g *Game) VotedOut() player.Player {
	// if an even number, first to half. otherwise, 50%+1
	threshold := len(g.PlayerList) / 2
	if len(g.PlayerList)%2 == 1 {
		threshold++
	}

	votes := 0
	var ousted player.Player
	for _, item := range g.Tally {
		votes += len(item.Votes)
		if len(item.Votes) >= threshold {
			ousted = item.Candidate
		}
	}
	if votes == len(g.PlayerList) && ousted != nil {
		return ousted
	}
	return nil
}

func (g *Game) QueueNightAction(fp *FingerPoint) {
	for i, a := range g.NightActionQueue {
		if a.From.ID() == fp.From.ID() {
			log.Printf("%s: replacing night action", fp.From.Identifier())
			g.NightActionQueue = append(g.NightActionQueue[:i], g.NightActionQueue[i+1:]...)
			break
		}
	}
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
	var death *player.Revealed
	for _, action := range g.NightActionQueue {
		result := g.NightAction(action)
		if result.PlayerMessage != "" {
			action.From.Message(message.Private, result.PlayerMessage)
		}

		// TODO maybe more than one death?
		if result.Killed != nil && death == nil && !result.Killed.Role().Kill() {
			death = result.Killed.Reveal()
			result.Killed.Message(message.Dead, "")
		}
	}

	g.NextPhase()

	if death != nil {
		g.Broadcast(message.Targeted, death)
	}
}

func (g *Game) Day() bool {
	return g.Phase%2 == 1
}

func (g *Game) RemovePlayer(id string) {
	delete(g.Players, id)
	g.UpdatePlayerList()

	// let everyone know they lost a friend
	g.Broadcast(message.PlayerList, g.PlayerList)

	if len(g.PlayerList) == 0 {
		g.Reset()
	}
}

func (g *Game) Reset() {
	for id, p := range g.Players {
		p.LeaveGame()
		delete(g.Players, id)
	}
	g.PlayerList = []player.Player{}
	g.Roleset = nil
	g.votes = make(map[player.Player]string)
	g.Tally = []*tally.TallyItem{}
	g.State = NotRunning
	g.Phase = 0
	g.NightActionQueue = []*FingerPoint{}
}

func (g *Game) Broadcast(title string, payload interface{}) {
	for _, p := range g.Players {
		if err := p.Message(title, payload); err != nil {
			log.Printf("%s: error messaging (%s)", p.ID(), err)
		}
	}
}

func (g *Game) HandleJoins(joinChan chan player.Player) {
	for {
		p := <-joinChan
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

		case chanmsg.GetRolesets:
			log.Printf("%s: requesting available rolesets", from.Identifier())
			from.Message(message.RolesetList, roleset.List())

		case chanmsg.SetRoleset:
			log.Printf("%s: setting roleset %s", from.Identifier(), activity.Roleset)
			sets := roleset.List()
			if err := g.SetRoleset(sets[activity.Roleset]); err != nil {
				log.Printf("error: %s", err)
				from.Message(message.Error, err)
			}

		case chanmsg.Vote:
			log.Printf("%s: voting for %s", from.Identifier(), g.Players[activity.To].Identifier())
			if err := g.Vote(g.Players[activity.From], activity.To); err != nil {
				// TODO?
			}

		case chanmsg.NightAction:
			log.Printf("%s: night action submitted", from.Identifier())
			fp := &FingerPoint{
				From: from,
				To:   to,
			}
			g.QueueNightAction(fp)

		case chanmsg.ResetGame:
			log.Printf("%s: resetting game", from.Identifier())
			g.Reset()
		}
	}
}
