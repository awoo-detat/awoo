package chanmsg

import (
	"github.com/gofrs/uuid"
)

const (
	Join       = iota
	SetName    = iota
	SetRoleset = iota
	PlayerList = iota
	Vote       = iota
	Quit       = iota
)

type Activity struct {
	Type int
	From *uuid.UUID
	To   *uuid.UUID
	//Roleset
}

func New(Type int, From *uuid.UUID) *Activity {
	return &Activity{
		Type: Type,
		From: From,
	}
}
