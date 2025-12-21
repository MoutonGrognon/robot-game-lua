package db

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/MoutonGrognon/robot-game-lua/matchmaker/internal/domain/entities"
	"github.com/MoutonGrognon/robot-game-lua/matchmaker/internal/domain/repositories"
)

type MatchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) repositories.MatchRepository {
	return &MatchRepository{
		db: db,
	}
}

func (mr *MatchRepository) Save(match entities.Match) error {
	stmt, err := mr.db.Prepare("INSERT INTO matchs (id, botId1, botId2, botName1, botName2, userName1, userName2, date, compressedGame, score1, score2, ranked) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(match.Id, match.BotId1, match.BotId2, match.BotName1, match.BotName2, match.UserName1, match.UserName2, match.Date, match.CompressedGame, match.Score1, match.Score2, match.Ranked)
	return err
}
