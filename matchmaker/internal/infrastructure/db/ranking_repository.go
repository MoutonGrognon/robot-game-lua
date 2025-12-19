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
	stmt, err := br.db.Prepare("SELECT id, botId, botName, elo, winCount, drawCount, lossCount FROM ranking")
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
		err = rows.Scan(&rank.Id, &rank.BotId, &rank.BotName, &rank.Elo, &rank.WinCount, &rank.DrawCount, &rank.LossCount)
		if err != nil {
			return ranks, err
		}
		ranks = append(ranks, rank)
	}
	err = rows.Err()
	return ranks, err
}

func (br *RankingRepository) UpdateRank(rank entities.Rank) error {
	stmt, err := br.db.Prepare("INSERT INTO ranking (id, botId, botName, elo, winCount, drawCount, lossCount) VALUES($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (id) DO UPDATE SET elo = EXCLUDED.elo, winCount = EXCLUDED.winCount, drawCount = EXCLUDED.drawCount, lossCount = EXCLUDED.lossCount;")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(rank.Id, rank.BotId, rank.BotName, rank.Elo, rank.WinCount, rank.DrawCount, rank.LossCount)
	return err
}
