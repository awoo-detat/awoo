package roleset

import (
	"stash.corp.synacor.com/hack/werewolf/role"
)

func VanillaFiver() *Roleset {
	return &Roleset{
		Name:        "Vanilla Fiver",
		Description: "Four villagers. One wolf. Two days to find them.",
		Roles: []*role.Role{
			role.Werewolf(),
			role.Villager(),
			role.Villager(),
			role.Villager(),
			role.Villager(),
		},
	}
}
