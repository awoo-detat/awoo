package player

import (
	"github.com/awoo-detat/awoo/chanmsg"

	"github.com/gofrs/uuid"
)

type PlayerType interface {
	Identifier() string
	Reveal() *Revealed
	Message(title string, payload interface{}) error
	SetChan(c chan *chanmsg.Activity)
	Vote(to uuid.UUID)
	NightAction(to uuid.UUID)
	Play()
	Reconnect(c Communicator)
	Quit()
}
