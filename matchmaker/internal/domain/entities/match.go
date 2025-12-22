package entities

import (
	"time"

	"github.com/google/uuid"

	"github.com/MoutonGrognon/robot-game-lua/rgcore/rgentities"
)

type Game []map[int]rgentities.BotState

type Match struct {
	Id             uuid.UUID `json:"id"`
	BlueBotId      uuid.UUID `json:"blueBotId"`
	RedBotId       uuid.UUID `json:"redBotId"`
	BlueBotName    string    `json:"blueBotName"`
	RedBotName     string    `json:"redBotName"`
	BlueUserName   string    `json:"blueUserName"`
	RedUserName    string    `json:"redUserName"`
	Date           time.Time `json:"date"`
	BlueScore      int       `json:"blueScore"`
	RedScore       int       `json:"redScore"`
	Ranked         bool      `json:"ranked"`
	CompressedGame []byte    `json:"compressedGame"`
}
