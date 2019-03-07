package role

func Cultist() *Role {
	return &Role{
		Name:           "Cultist",
		Description:    "You know the wolf, but they don't know you. Lie your way to an evil win!",
		Team:           Evil,
		VoteMultiplier: 1,
		Health:         1,
		Parity:         1,
		Alive:          true,
		Attributes:     auxEvil,
		Actions:        knowsMaxes,
	}
}
