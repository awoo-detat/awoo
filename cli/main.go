package main

import (
	"fmt"

	"stash.corp.synacor.com/hack/werewolf/cli/communicator"
	"stash.corp.synacor.com/hack/werewolf/game"
	"stash.corp.synacor.com/hack/werewolf/player"
	"stash.corp.synacor.com/hack/werewolf/role/roleset"
)

func main() {
	c := make(chan *player.Player)
	game := game.New(c)
	game.SetRoleset(roleset.Fiver())

	dan := player.New(communicator.New("dan.log"), c)
	dan.Name = "Dan"
	game.AddPlayer(dan)
	jon := player.New(communicator.New("jon.log"), c)
	game.AddPlayer(jon)
	tyler := player.New(communicator.New("tyler.log"), c)
	game.AddPlayer(tyler)
	matt := player.New(communicator.New("matt.log"), c)
	game.AddPlayer(matt)
	matt.Name = "Matt"
	julia := player.New(communicator.New("julia.log"), c)
	julia.Name = "Julia"
	game.AddPlayer(julia)

}
