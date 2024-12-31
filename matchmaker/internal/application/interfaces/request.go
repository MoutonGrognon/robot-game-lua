package interfaces

import (
	"github.com/MoutonGrognon/robot-game-lua/rgcore"
)

type SaveMatchRequest struct {
	MatchId string                    `json:"matchId"`
	Match   []map[int]rgcore.BotState `json:"match"`
}

type CancelMatchRequest struct {
	MatchId string `json:"matchId"`
	Error   error  `json:"error"`
}
