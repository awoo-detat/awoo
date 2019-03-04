package roleset

import (
	"github.com/awoo-detat/awoo/role"
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

func init() {
	registerRoleset(VanillaFiver())
}
