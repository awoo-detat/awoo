package main

import (
	"log"
	"time"
	//"fmt"

	"stash.corp.synacor.com/hack/werewolf/cli/communicator"
	"stash.corp.synacor.com/hack/werewolf/game"
	"stash.corp.synacor.com/hack/werewolf/player"
	"stash.corp.synacor.com/hack/werewolf/role/roleset"
)

func main() {
	for _, set := range roleset.List() {
		log.Printf("%s: %v players", set.Name, len(set.Roles))
	}
	c := make(chan *player.Player)
	game := game.New(c)

	dan := player.New(communicator.New("dan.log"), c)
	dan.Name = "Dan"
	game.AddPlayer(dan)
	game.SetRoleset(roleset.Fiver())
	jon := player.New(communicator.New("jon.log"), c)
	jon.Name = "Jon"
	game.AddPlayer(jon)
	tyler := player.New(communicator.New("tyler.log"), c)
	tyler.Name = "Tyler"
	game.AddPlayer(tyler)
	matt := player.New(communicator.New("matt.log"), c)
	game.AddPlayer(matt)
	matt.Name = "Matt"
	julia := player.New(communicator.New("julia.log"), c)
	julia.Name = "Julia"
	game.AddPlayer(julia)

	// game has five players and has started

	var wolf *player.Player
	var villagers []*player.Player

	// hack the planet
	if dan.Role.Name == "Werewolf" {
		wolf = dan
	} else {
		villagers = append(villagers, dan)
	}
	if jon.Role.Name == "Werewolf" {
		wolf = jon
	} else {
		villagers = append(villagers, jon)
	}
	if tyler.Role.Name == "Werewolf" {
		wolf = tyler
	} else {
		villagers = append(villagers, tyler)
	}
	if matt.Role.Name == "Werewolf" {
		wolf = matt
	} else {
		villagers = append(villagers, matt)
	}
	if julia.Role.Name == "Werewolf" {
		wolf = julia
	} else {
		villagers = append(villagers, julia)
	}

	wolf.Vote(villagers[0].UUID)
	villagers[1].Vote(villagers[0].UUID)
	villagers[2].Vote(villagers[0].UUID)

	villagers[1].NightAction(villagers[2].UUID)
	villagers[2].NightAction(villagers[1].UUID)
	villagers[3].NightAction(villagers[2].UUID)
	wolf.NightAction(villagers[1].UUID)

	wolf.Vote(villagers[2].UUID)
	villagers[2].Vote(wolf.UUID)
	villagers[3].Vote(wolf.UUID)
	//villagers[3].Vote(villagers[2].UUID)

	time.Sleep(1 * time.Second)
}
