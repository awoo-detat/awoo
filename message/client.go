package message

import (
	"encoding/json"
	"log"
)

type ClientMessage struct {
	JoinName       string `json:"joinWithName"`
	PollPlayerList bool   `json:"pollPlayerList"`
}

func Decode(raw []byte) *ClientMessage {
	var message ClientMessage
	err := json.Unmarshal(raw, &message)
	if err != nil {
		log.Printf("error decoding client message: %s", err)
	}
	return &message
}
