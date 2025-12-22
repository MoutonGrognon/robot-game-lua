package entities

import (
	"github.com/google/uuid"
)

type PendingMatch struct {
	Id        uuid.UUID `json:"id"`
	BotId1    uuid.UUID `json:"botId1"`
	BotId2    uuid.UUID `json:"botId2"`
	BotName1  string    `json:"botName1"`
	BotName2  string    `json:"botName2"`
	UserName1 string    `json:"userName1"`
	UserName2 string    `json:"userName2"`
	Ranked    bool      `json:"ranked"`
}
