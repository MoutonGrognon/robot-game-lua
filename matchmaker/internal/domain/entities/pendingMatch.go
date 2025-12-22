package entities

import (
	"github.com/google/uuid"
)

type PendingMatch struct {
	Id           uuid.UUID `json:"id"`
	BlueBotId    uuid.UUID `json:"blueBotId"`
	RedBotId     uuid.UUID `json:"redBotId"`
	BlueBotName  string    `json:"blueBotName"`
	RedBotName   string    `json:"redBotName"`
	BlueUserName string    `json:"blueUserName"`
	RedUserName  string    `json:"redUserName"`
	Ranked       bool      `json:"ranked"`
}
