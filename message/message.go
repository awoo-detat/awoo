package message

import (
	"encoding/json"
)

const (
	Awoo       = "awoo"
	Role       = "role"
	PlayerList = "playerlist"
	Roleset    = "roleset"
	GameState  = "gamestate"
	Phase      = "phase"
	Leader     = "leader"
	PleaseWait = "pleasewait" // they're in the lobby
	CanNotJoin = "cannotjoin"
)

type Message struct {
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

func New(message string, payload interface{}) ([]byte, error) {
	m := Message{
		Message: message,
		Payload: payload,
	}

	b, err := json.Marshal(m)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}
