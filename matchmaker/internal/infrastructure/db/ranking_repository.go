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

func (rr *RankingRepository) GetRanking() ([]entities.Rank, error) {
	var ranks []entities.Rank
	stmt, err := rr.db.Prepare("SELECT id, botId, botName, elo, winCount, drawCount, lossCount, streak FROM ranking")
	if err != nil {
		return ranks, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return ranks, err
	}
	defer rows.Close()
	for rows.Next() {
		var rank entities.Rank
		err = rows.Scan(&rank.Id, &rank.BotId, &rank.BotName, &rank.Elo, &rank.WinCount, &rank.DrawCount, &rank.LossCount, &rank.Streak)
		if err != nil {
			return ranks, err
		}
		ranks = append(ranks, rank)
	}
	err = rows.Err()
	return ranks, err
}

func (rr *RankingRepository) UpdateRank(rank entities.Rank) error {
	stmt, err := rr.db.Prepare("INSERT INTO ranking (id, botId, botName, elo, winCount, drawCount, lossCount, streak) VALUES($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (id) DO UPDATE SET elo = EXCLUDED.elo, winCount = EXCLUDED.winCount, drawCount = EXCLUDED.drawCount, lossCount = EXCLUDED.lossCount, streak = EXCLUDED.streak;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(rank.Id, rank.BotId, rank.BotName, rank.Elo, rank.WinCount, rank.DrawCount, rank.LossCount, rank.Streak)
	return err
}
