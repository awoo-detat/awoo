package role

func Hunter() *Role {
	return &Role{
		Name:           "Hunter",
		Description:    "Evil can't win unless you're dead. So... don't do that.",
		Team:           Good,
		VoteMultiplier: 1,
		Health:         1,
		Parity:         2,
		Alive:          true,
	}
}
