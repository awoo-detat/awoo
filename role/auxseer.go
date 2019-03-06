package role

func AuxSeer() *Role {
	return &Role{
		Name:           "Aux Seer",
		Description:    "Each night you choose someone to view, and are told if they're aux evil (a Cultist, Sorcerer, etc). At the start of the game you're given a random player who is not aux evil.",
		Team:           Good,
		VoteMultiplier: 1,
		Health:         1,
		Parity:         1,
		Alive:          true,
		Actions:        viewForAux | randomN0Clear,
		Attributes:     seer,
	}
}
