package repositories

import (
	"github.com/MoutonGrognon/robot-game-lua/referee/internal/domain/entities"
)

type BotRepository interface {
	GetByName(name string) (entities.Bot, error)
}
