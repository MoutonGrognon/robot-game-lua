package entities

import "github.com/google/uuid"

type Rank struct {
	Id        uuid.UUID `json:"id"`
	BotId     uuid.UUID `json:"botId"`
	BotName   string    `json:"botName"`
	Elo       int       `json:"elo"`
	WinCount  int       `json:"winCount"`
	DrawCount int       `json:"drawCount"`
	LossCount int       `json:"lossCount"`
}
