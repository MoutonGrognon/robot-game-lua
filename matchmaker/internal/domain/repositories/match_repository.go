package repositories

import (
	"github.com/MoutonGrognon/robot-game-lua/matchmaker/internal/domain/entities"
)

type MatchRepository interface {
	Save(match entities.Match) error
}
