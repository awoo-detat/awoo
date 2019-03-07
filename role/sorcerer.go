package role

func Sorcerer() *Role {
	return &Role{
		Name:           "Sorcerer",
		Description:    "You don't know the wolf, but are evil. Each night you view a player and find out if they're the seer. Kill that seer.",
		Team:           Evil,
		VoteMultiplier: 1,
		Health:         1,
		Parity:         1,
		Alive:          true,
		Actions:        randomN0Clear | viewForSeer,
		Attributes:     auxEvil,
	}
}
