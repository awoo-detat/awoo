package role

func Werewolf() *Role {
	return &Role{
		Name:           "Werewolf",
		Description:    "Blend in during the day. At night... you feed.",
		Team:           Evil,
		VoteMultiplier: 1,
		Health:         1,
		Parity:         -1,
		Alive:          true,
		Actions:        knowsMaxes | nightKill,
		Attributes:     maxEvil,
	}
}
