package entities

import "github.com/google/uuid"

type Rank struct {
	Id      uuid.UUID
	BotId   uuid.UUID
	BotName string
	Elo     int
}
