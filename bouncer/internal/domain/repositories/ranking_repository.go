package repositories

import "github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/entities"

type RankingRepository interface {
	GetRanking() ([]entities.Rank, error)
}
