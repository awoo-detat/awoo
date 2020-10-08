package game

import (
	"fmt"
	"log"
	"strings"

	"github.com/awoo-detat/awoo/player"
)

const (
	IsWerewolf       = "IS a Werewolf!"
	IsNotWerewolf    = "is not a Werewolf"
	IsSeer           = "IS a Seer"
	IsNotSeer        = "is not a Seer"
	IsAuxEvil        = "IS a aux evil"
	IsNotAuxEvil     = "is not aux evil"
	WolfListSingle   = "The Werewolf is"
	WolfListMultiple = "The Werewolves are"
)

type ActionResult struct {
	PlayerMessage string
	Killed        player.Player
}

func (g *Game) NightAction(fp *FingerPoint) *ActionResult {
	result := &ActionResult{}

	if fp.From.Role().HasNightKill() {
		log.Printf("%s is killing %s", fp.From.Identifier(), fp.To.Identifier())
		result.Killed = fp.To
	}

	if fp.From.Role().ViewsForMax() {
		log.Printf("%s (seer) is viewing %s", fp.From.Identifier(), fp.To.Identifier())
		if fp.To.Role().ViewForMaxEvil() {
			result.PlayerMessage = fmt.Sprintf("%s %s", fp.To.Identifier(), IsWerewolf)
		} else {
			result.PlayerMessage = fmt.Sprintf("%s %s", fp.To.Identifier(), IsNotWerewolf)
		}
	}

	if fp.From.Role().ViewsForSeer() {
		log.Printf("%s (sorcerer) is viewing %s", fp.From.Identifier(), fp.To.Identifier())
		if fp.To.Role().ViewForSeer() {
			result.PlayerMessage = fmt.Sprintf("%s %s", fp.To.Identifier(), IsSeer)
		} else {
			result.PlayerMessage = fmt.Sprintf("%s %s", fp.To.Identifier(), IsNotSeer)
		}
	}

	if fp.From.Role().ViewsForAux() {
		log.Printf("%s (aux seer) is viewing %s", fp.From.Identifier(), fp.To.Identifier())
		if fp.To.Role().ViewForAuxEvil() {
			result.PlayerMessage = fmt.Sprintf("%s %s", fp.To.Identifier(), IsAuxEvil)
		} else {
			result.PlayerMessage = fmt.Sprintf("%s %s", fp.To.Identifier(), IsNotAuxEvil)
		}
	}

	return result
}

func (g *Game) StartAction(p player.Player) *ActionResult {
	result := &ActionResult{}

	if p.Role().KnowsMaxes() {
		maxes := g.AliveMaxEvils()
		if len(maxes) > 1 {
			result.PlayerMessage = fmt.Sprintf("%s %s", WolfListMultiple, strings.Join(maxes, ", "))
		} else if len(maxes) == 1 && !p.Role().IsMaxEvil() {
			// don't send the Wolf a PM if they're the only one
			result.PlayerMessage = fmt.Sprintf("%s %s", WolfListSingle, maxes[0])
		}
	}
	if p.Role().HasRandomN0Clear() {
		var clear string
		var role string
		// this, uh, isn't really random. TODO?
		for _, player := range g.Players {
			var hit bool
			switch {
			case p.Role().ViewsForMax():
				hit = player.Role().ViewForMaxEvil()
				role = IsNotWerewolf
			case p.Role().ViewsForSeer():
				hit = player.Role().ViewForSeer()
				role = IsNotSeer
			case p.Role().ViewsForAux():
				hit = player.Role().ViewForAuxEvil()
				role = IsNotAuxEvil
			}

			if p.ID() != player.ID() && !hit {
				clear = player.Identifier()
				break
			}
		}
		if clear != "" {
			result.PlayerMessage = fmt.Sprintf("%s %s", clear, role)
		}
	}

	return result
}
