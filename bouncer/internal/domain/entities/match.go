package entities

import (
	"time"

	"github.com/google/uuid"
)

type MatchSummary struct {
	Id           uuid.UUID `json:"id"`
	BlueBotId    uuid.UUID `json:"blueBotId"`
	RedBotId     uuid.UUID `json:"redBotId"`
	BlueBotName  string    `json:"blueBotName"`
	RedBotName   string    `json:"redBotName"`
	BlueUserName string    `json:"blueUserName"`
	RedUserName  string    `json:"redUserName"`
	Date         time.Time `json:"date"`
	BlueScore    int       `json:"blueScore"`
	RedScore     int       `json:"redScore"`
	Ranked       bool      `json:"ranked"`
}

type Match struct {
	MatchSummary
	CompressedGame []byte `json:"compressedGame"`
}
