package roleset

import (
	"stash.corp.synacor.com/hack/werewolf/role"
)

func Niner() *Roleset {
	return &Roleset{
		Name:        "Basic Niner",
		Description: "There are two wolves to find, so keep your seer protected!",
		Roles: []*role.Role{
			role.Werewolf(),
			role.Werewolf(),
			role.Cultist(),
			role.Hunter(),
			role.Seer(),
			role.Villager(),
			role.Villager(),
			role.Villager(),
			role.Villager(),
		},
	}
}

func init() {
	registerRoleset(Niner())
}
