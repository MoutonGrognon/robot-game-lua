package interfaces

import (
	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/entities"
)

type GetMatchResponse struct {
	Match entities.Match `json:"match"`
}

type GetSummariesResponse struct {
	Summaries []entities.MatchSummary `json:"summaries"`
	Start     int                     `json:"start"`
	Size      int                     `json:"size"`
	Total     int                     `json:"total"`
}
