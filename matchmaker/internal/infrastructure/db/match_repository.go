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
	stmt, err := mr.db.Prepare("INSERT INTO matches (id, blueBotId, redBotId, blueBotName, redBotName, blueUserName, redUserName, date, compressedGame, blueScore, redScore, ranked) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(match.Id, match.BlueBotId, match.RedBotId, match.BlueBotName, match.RedBotName, match.BlueUserName, match.RedUserName, match.Date, match.CompressedGame, match.BlueScore, match.RedScore, match.Ranked)
	return err
}
