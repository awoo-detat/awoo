package chanmsg

import (
	"github.com/gofrs/uuid"
)

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
	From    uuid.UUID
	To      uuid.UUID
	Roleset string
}

func New(Type int, From uuid.UUID) *Activity {
	return &Activity{
		Type: Type,
		From: From,
	}
}
