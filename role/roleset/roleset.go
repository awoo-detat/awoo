package roleset

import (
	"stash.corp.synacor.com/hack/werewolf/role"
)

type Roleset struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Roles       []*role.Role `json:"roles"`
}
