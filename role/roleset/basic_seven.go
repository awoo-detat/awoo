package roleset

import (
	"github.com/awoo-detat/awoo/role"
)

func BasicSeven() *Roleset {
	return &Roleset{
		Name:        "Basic Seven",
		Description: "One wolf, one cultist",
		Roles: []*role.Role{
			role.Werewolf(),
			role.Cultist(),
			role.Villager(),
			role.Villager(),
			role.Villager(),
			role.Villager(),
			role.Villager(),
		},
	}
}

func init() {
	registerRoleset(BasicSeven())
}
