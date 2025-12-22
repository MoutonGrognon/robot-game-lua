package entities

import "github.com/google/uuid"

type Bot struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Script string    `json:"script"`
}
