package roleset

import (
	"github.com/awoo-detat/awoo/role"
)

func LessBasicSeven() *Roleset {
	return &Roleset{
		Name:        "Less Basic Seven",
		Description: "Seer and sorcerer, oh my!",
		Roles: []*role.Role{
			role.Werewolf(),
			role.Sorcerer(),
			role.Seer(),
			role.Villager(),
			role.Villager(),
			role.Villager(),
			role.Villager(),
		},
	}
}

func init() {
	registerRoleset(LessBasicSeven())
}
