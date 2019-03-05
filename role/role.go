package role

const (
	Good    = iota
	Evil    = iota
	Neutral = iota // won't be used during the Hackathon
)

type Attribute int

const (
	maxEvil Attribute = 1 << iota
	auxEvil
	seer
	tinker
)

type Action int

const (
	viewForMax Action = 1 << iota
	nightKill
	viewForSeer
	viewForAux
	randomN0Clear
	knowsMaxes
)

type Role struct {
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Team           int       `json:"team"`
	Parity         int       `json:"-"`
	VoteMultiplier int       `json:"-"`
	Health         int       `json:"-"`
	Alive          bool      `json:"alive"`
	Actions        Action    `json:"night_action"`
	Attributes     Attribute `json:"-"`
}

// IsMaxEvil returns whether or not a player is a max evil (ie a Werewolf).
func (r *Role) IsMaxEvil() bool {
	return r.Attributes&maxEvil > 0
}

// ViewForMaxEvil allows seers to view if a role is max evil. It differs from IsMaxEvil because the
// Tinker can invert the result.
func (r *Role) ViewForMaxEvil() bool {
	if r.Attributes&tinker > 0 {
		return !r.IsMaxEvil()
	}
	return r.IsMaxEvil()
}

// IsAuxEvil returns whether or not a player is an aux evil (ie a Cultist).
func (r *Role) IsAuxEvil() bool {
	return r.Attributes&auxEvil > 0
}

// ViewForAuxEvil allows seers to view if a role is aux evil. It differs from IsAuxEvil because the Tinker
// can invert the result.
func (r *Role) ViewForAuxEvil() bool {
	if r.Attributes&tinker > 0 {
		return !r.IsAuxEvil()
	}
	return r.IsAuxEvil()
}

// IsSeer returns whether or not a player is a seer
func (r *Role) IsSeer() bool {
	return r.Attributes&seer > 0
}

// ViewForSeer allows sorcerers to view if a role is a seer. It differs from IsSeer because the
// Tinker can invert the result.
func (r *Role) ViewForSeer() bool {
	if r.Attributes&tinker > 0 {
		return !r.IsSeer()
	}
	return r.IsSeer()
}

func (r *Role) ViewsForMax() bool {
	return r.Actions&viewForMax > 0
}

func (r *Role) HasNightKill() bool {
	return r.Actions&nightKill > 0
}

func (r *Role) ViewsForSeer() bool {
	return r.Actions&viewForSeer > 0
}

func (r *Role) ViewsForAux() bool {
	return r.Actions&viewForAux > 0
}

func (r *Role) HasRandomN0Clear() bool {
	return r.Actions&randomN0Clear > 0
}

func (r *Role) KnowsMaxes() bool {
	return r.Actions&knowsMaxes > 0
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
