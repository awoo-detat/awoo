package role

func Villager() *Role {
	return &Role{
		Name:           "Villager",
		Description:    "You have no special powers. Kill the wolves before they kill you!",
		Team:           Good,
		VoteMultiplier: 1,
		Health:         1,
		Parity:         1,
		Alive:          true,
	}
}
