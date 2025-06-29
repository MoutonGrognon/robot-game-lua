package repositories

import (
<<<<<<< HEAD
	"github.com/MoutonGrognon/robot-game-lua/matchmaker/internal/domain/entities"
=======
	"github.com/MoutonGrognon/robot-game-lua/matchmaker/internal/domain/entities"
	"github.com/google/uuid"
>>>>>>> d97a808 (feat: WIP add frontend + add highlighted match endpoint)
)

type MatchRepository interface {
	Save(match entities.Match) error
	// TODO: WIP
	GetById(id uuid.UUID) (entities.Match, error)
}
