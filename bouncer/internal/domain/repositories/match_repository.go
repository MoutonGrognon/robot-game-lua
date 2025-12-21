package repositories

import (
	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/entities"
	"github.com/google/uuid"
)

type MatchRepository interface {
	GetById(id uuid.UUID) (entities.Match, error)
	GetSummaries(start int, size int) ([]entities.MatchSummary, int, int, int, error)
	// TODO: WIP
}
