package repositories

import "github.com/MoutonGrognon/robot-game-lua/matchmaker/internal/domain/entities"

type RankingRepository interface {
	GetRanking() ([]entities.Rank, error)
}
