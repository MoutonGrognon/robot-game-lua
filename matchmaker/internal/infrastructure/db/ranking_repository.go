package db

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/MoutonGrognon/robot-game-lua/matchmaker/internal/domain/entities"
	"github.com/MoutonGrognon/robot-game-lua/matchmaker/internal/domain/repositories"
)

type RankingRepository struct {
	db *sql.DB
}

func NewRankingRepository(db *sql.DB) repositories.RankingRepository {
	return &RankingRepository{
		db: db,
	}
}

func (br *RankingRepository) GetRanking() ([]entities.Rank, error) {
	var ranks []entities.Rank
	stmt, err := br.db.Prepare("SELECT id, botId, botName, elo FROM ranking")
	if err != nil {
		return ranks, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return ranks, err
	}
	defer rows.Close()
	for rows.Next() {
		var rank entities.Rank
		err = rows.Scan(&rank.Id, &rank.BotId, &rank.BotName, &rank.Elo)
		if err != nil {
			return ranks, err
		}
		ranks = append(ranks, rank)
	}
	err = rows.Err()
	return ranks, err
}
