package player

import (
	"github.com/awoo-detat/awoo/chanmsg"
	"github.com/awoo-detat/awoo/role"

	"github.com/gofrs/uuid"
)

type Player interface {
	UUID() uuid.UUID
	Identifier() string
	Reveal() *Revealed
	Message(title string, payload interface{}) error
	SetChan(c chan *chanmsg.Activity)
	Vote(to uuid.UUID)
	NightAction(to uuid.UUID)
	Play()
	InGame() bool
	LeaveGame()
	SetLeader()
	SetName(name string)
	SetRole(r *role.Role)
	Reconnect(c Communicator)
	Role() *role.Role
	Quit()
}
