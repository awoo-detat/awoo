package chanmsg

const (
	Join = iota
	SetName
	SetRoleset
	PlayerList
	Vote
	Quit
	Tally
	NightAction
	ResetGame
)

type Activity struct {
	Type    int
	From    string
	To      string
	Roleset string
}

func New(Type int, From string) *Activity {
	return &Activity{
		Type: Type,
		From: From,
	}
}
