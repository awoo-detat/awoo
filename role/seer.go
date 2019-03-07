package role

func Seer() *Role {
	return &Role{
		Name:           "Seer",
		Description:    "Each night you choose someone to view, and are told if they're a Werewolf. At the start of the game you're given a random player who is not a Wolf.",
		Team:           Good,
		VoteMultiplier: 1,
		Health:         1,
		Parity:         1,
		Alive:          true,
		Actions:        viewForMax | randomN0Clear,
		Attributes:     seer,
	}
}
