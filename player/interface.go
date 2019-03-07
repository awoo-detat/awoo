package player

import (
	"github.com/awoo-detat/awoo/chanmsg"
	"github.com/awoo-detat/awoo/role"
)

type Player interface {
	ID() string
	Identifier() string
	Reveal() *Revealed
	Message(title string, payload interface{}) error
	SetChan(c chan *chanmsg.Activity)
	Vote(to string)
	NightAction(to string)
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
