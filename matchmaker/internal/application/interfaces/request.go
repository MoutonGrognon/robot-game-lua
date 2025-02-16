package interfaces

import (
	"github.com/google/uuid"

	"github.com/MoutonGrognon/robot-game-lua/rgcore"
)

type SaveMatchRequest struct {
	MatchId uuid.UUID                 `json:"matchId"`
	Match   []map[int]rgcore.BotState `json:"match"`
}

type CancelMatchRequest struct {
	MatchId uuid.UUID `json:"matchId"`
	Error   error     `json:"error"`
}
