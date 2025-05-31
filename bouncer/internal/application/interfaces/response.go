package interfaces

import (
	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/entities"
)

type GetMatchResponse struct {
	Match entities.Match `json:"match"`
}
