package role

const (
	Good    = iota
	Evil    = iota
	Neutral = iota // won't be used during the Hackathon
)

type Role struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	Team           int    `json:"team"`
	Parity         int    `json:"-"`
	VoteMultiplier int    `json:"-"`
	Health         int    `json:"-"`
	Alive          bool   `json:"alive"`
	NightAction    bool   `json:"night_action"`
}

// Kill attempts to kill the player. If they had more than 1 health (ie
// were "tough") then they will remain alive.
func (r *Role) Kill() bool {
	// maybe this should error if you try to kill a dead person?
	r.Health--
	if r.Health <= 0 {
		r.Alive = false
	}
	return r.Alive
}
