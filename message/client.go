package message

import (
	"encoding/json"
	"log"
)

const (
	TimeDay   = "day"
	TimeNight = "night"
)

type ClientMessage struct {
	JoinName       string `json:"joinWithName"`
	PollPlayerList bool   `json:"pollPlayerList"`
	PollTally      bool   `json:"pollTally"`
	Vote           string `json:"voteFor"`
	Time           string `json:"time"`
}

func Decode(raw []byte) *ClientMessage {
	var message ClientMessage
	err := json.Unmarshal(raw, &message)
	if err != nil {
		log.Printf("error decoding client message: %s", err)
	}
	return &message
}
