package repositories

import (
	"time"

	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/entities"
	"github.com/google/uuid"
)

type MatchRepository interface {
	GetById(id uuid.UUID) (entities.Match, error)
	GetSummaries(start int, size int) ([]entities.MatchSummary, int, int, int, error)
	GetRecentSummaries(dateThreshold time.Time) ([]entities.MatchSummary, error)
	DeleteOldMatches(int) (int, error)
}
