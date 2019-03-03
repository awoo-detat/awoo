package roleset

import (
	"github.com/Sigafoos/awoo/role"
)

type Roleset struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Roles       []*role.Role `json:"roles"`
}

func List() map[string]*Roleset {
	return sets
}

type RolesetMap map[string]*Roleset

var sets = RolesetMap{}

func registerRoleset(roleset *Roleset) {
	sets[roleset.Name] = roleset
}
