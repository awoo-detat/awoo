package roleset

import (
	"github.com/awoo-detat/awoo/role"
)

func ImperfectEleven() *Roleset {
	return &Roleset{
		Name:        "Imperfect Eleven",
		Description: "Two wolves and a semblance of balance.",
		Roles: []*role.Role{
			role.Werewolf(),
			role.Werewolf(),
			role.Sorcerer(),
			role.Hunter(),
			role.Seer(),
			role.Villager(),
			role.Villager(),
			role.Villager(),
			role.Villager(),
			role.Villager(),
			role.Villager(),
		},
	}
}

func init() {
	registerRoleset(ImperfectEleven())
}
