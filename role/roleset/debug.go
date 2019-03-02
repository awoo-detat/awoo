package roleset

import (
	"stash.corp.synacor.com/hack/werewolf/role"
)

func Debug() *Roleset {
	return &Roleset{
		Name:        "awooooooooo",
		Description: "to save jcantwell's sanity",
		Roles: []*role.Role{
			role.Werewolf(),
			role.Hunter(),
			role.Villager(),
		},
	}
}
