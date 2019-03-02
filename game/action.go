package game

import (
	"fmt"
	"log"
	"strings"

	"stash.corp.synacor.com/hack/werewolf/player"
	"stash.corp.synacor.com/hack/werewolf/role"
)

type ActionResult struct {
	PlayerMessage string
	Killed        *player.Player
}

func (g *Game) NightAction(fp *FingerPoint) *ActionResult {
	result := &ActionResult{}

	if fp.From.Role.NightAction&role.NightKill > 0 {
		log.Printf("%s is killing %s", fp.From.Identifier(), fp.To.Identifier())
		if !fp.To.Role.Kill() {
			result.Killed = fp.To
		}
	}

	if fp.From.Role.NightAction&role.ViewForMax > 0 {
		log.Printf("%s (seer) is viewing %s", fp.From.Identifier(), fp.To.Identifier())
		if fp.To.Role.Parity < 0 {
			result.PlayerMessage = fmt.Sprintf("%s IS a Werewolf!", fp.To.Identifier())
		} else {
			result.PlayerMessage = fmt.Sprintf("%s is not a Werewolf", fp.To.Identifier())
		}
	}

	return result
}

func (g *Game) StartAction(p *player.Player) *ActionResult {
	result := &ActionResult{}

	if p.Role.StartAction&role.MaxList > 0 {
		maxes := g.AliveMaxEvils()
		if len(maxes) > 1 {
			result.PlayerMessage = fmt.Sprintf("The Werewolves are %s", strings.Join(maxes, ", "))
		} else if len(maxes) == 1 && p.Role.Parity > 0 {
			// don't send the Wolf a PM if they're the only one
			result.PlayerMessage = fmt.Sprintf("The Werewolf is %s", maxes[0])
		}
	}
	if p.Role.StartAction&role.RandomMaxClear > 0 {
		var clear string
		for _, player := range g.Players {
			if p.UUID != player.UUID && player.Role.Parity > 0 {
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
