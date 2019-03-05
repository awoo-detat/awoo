package game

import (
	"fmt"
	"log"
	"strings"

	"github.com/awoo-detat/awoo/player"
	"github.com/awoo-detat/awoo/role"
)

type ActionResult struct {
	PlayerMessage string
	Killed        *player.Player
}

func (g *Game) NightAction(fp *FingerPoint) *ActionResult {
	result := &ActionResult{}

	if fp.From.Role.Actions&role.NightKill > 0 {
		log.Printf("%s is killing %s", fp.From.Identifier(), fp.To.Identifier())
		if !fp.To.Role.Kill() {
			result.Killed = fp.To
		}
	}

	if fp.From.Role.Actions&role.ViewForMax > 0 {
		log.Printf("%s (seer) is viewing %s", fp.From.Identifier(), fp.To.Identifier())
		if fp.To.Role.ViewForMaxEvil() {
			result.PlayerMessage = fmt.Sprintf("%s IS a Werewolf!", fp.To.Identifier())
		} else {
			result.PlayerMessage = fmt.Sprintf("%s is not a Werewolf", fp.To.Identifier())
		}
	}

	return result
}

func (g *Game) StartAction(p *player.Player) *ActionResult {
	result := &ActionResult{}

	if p.Role.Actions&role.ToldMaxes > 0 {
		maxes := g.AliveMaxEvils()
		if len(maxes) > 1 {
			result.PlayerMessage = fmt.Sprintf("The Werewolves are %s", strings.Join(maxes, ", "))
		} else if len(maxes) == 1 && p.Role.IsMaxEvil() {
			// don't send the Wolf a PM if they're the only one
			result.PlayerMessage = fmt.Sprintf("The Werewolf is %s", maxes[0])
		}
	}
	if p.Role.Actions&role.RandomN0Clear > 0 {
		var clear string
		// this, uh, isn't really random. TODO?
		for _, player := range g.Players {
			var hit bool
			switch {
			case p.Role.Actions&role.ViewForMax > 0:
				hit = player.Role.ViewForMaxEvil()
			case p.Role.Actions&role.ViewForSeer > 0:
				hit = player.Role.ViewForSeer()
			case p.Role.Actions&role.ViewForAux > 0:
				hit = player.Role.ViewForAuxEvil()
			}

			if p.UUID != player.UUID && !hit {
				clear = player.Identifier()
				break
			}
		}
		if clear != "" {
			result.PlayerMessage = fmt.Sprintf("%s is not a Werewolf", clear)
		}
	}

	return result
}
