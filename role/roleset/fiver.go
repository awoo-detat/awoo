package roleset

import (
	"github.com/awoo-detat/awoo/role"
)

func Fiver() *Roleset {
	return &Roleset{
		Name:        "Fast Fiver",
		Description: "Two chances to catch the wolf, but both sides have some power.",
		Roles: []*role.Role{
			role.Werewolf(),
			role.Cultist(),
			role.Hunter(),
			role.Seer(),
			role.Villager(),
		},
	}
}

func init() {
	registerRoleset(Fiver())
}
