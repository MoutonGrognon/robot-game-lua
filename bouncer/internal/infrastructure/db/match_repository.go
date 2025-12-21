package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/entities"
	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/repositories"
)

type MatchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) repositories.MatchRepository {
	return &MatchRepository{
		db: db,
	}
}

func (mr *MatchRepository) GetById(id uuid.UUID) (entities.Match, error) {
	var match entities.Match
	stmt, err := mr.db.Prepare("SELECT * FROM matchs WHERE id=$1")
	if err != nil {
		return match, err
	}
	err = stmt.QueryRow(id).Scan(&match.Id, &match.BotId1, &match.BotId2, &match.BotName1, &match.BotName2, &match.UserName1, &match.UserName2, &match.Date, &match.CompressedGame, &match.Score1, &match.Score2, &match.Ranked)
	if err != nil {
		return match, err
	}
	return match, err
}

func (mr *MatchRepository) GetSummaries(start int, size int) ([]entities.MatchSummary, int, int, int, error) {
	var matchs []entities.MatchSummary
	total := 0
	stmt, err := mr.db.Prepare("SELECT id, botId1, botId2, botName1, botName2, userName1, userName2, date, score1, score2, ranked FROM (" +
		"SELECT ROW_NUMBER() OVER (ORDER BY date DESC) as rowNum, * FROM matchs" +
		") WHERE (rowNum>=$1 AND rowNum<$2) ORDER BY rowNum")
	if err != nil {
		return matchs, start, size, total, err
	}
	rows, err := stmt.Query(start, start+size)
	if err != nil {
		return matchs, start, size, total, err
	}
	defer rows.Close()
	for rows.Next() {
		var match entities.MatchSummary
		err = rows.Scan(&match.Id, &match.BotId1, &match.BotId2, &match.BotName1, &match.BotName2, &match.UserName1, &match.UserName2, &match.Date, &match.Score1, &match.Score2, &match.Ranked)
		if err != nil {
			return matchs, start, size, total, err
		}
		matchs = append(matchs, match)
	}
	err = rows.Err()
	if err != nil {
		return matchs, start, size, total, err
	}
	row := mr.db.QueryRow("SELECT count FROM matchs_count")
	row.Scan(&total)
	return matchs, start, size, total, err
}

func (mr *MatchRepository) GetRecentSummaries(dateThreshold time.Time) ([]entities.MatchSummary, error) {
	var matchs []entities.MatchSummary
	stmt, err := mr.db.Prepare("SELECT id, botId1, botId2, botName1, botName2, userName1, userName2, date, score1, score2, ranked FROM matchs WHERE date>=$1")
	if err != nil {
		return matchs, err
	}
	rows, err := stmt.Query(dateThreshold)
	if err != nil {
		return matchs, err
	}
	defer rows.Close()
	for rows.Next() {
		var match entities.MatchSummary
		err = rows.Scan(&match.Id, &match.BotId1, &match.BotId2, &match.BotName1, &match.BotName2, &match.UserName1, &match.UserName2, &match.Date, &match.Score1, &match.Score2, &match.Ranked)
		if err != nil {
			return matchs, err
		}
		matchs = append(matchs, match)
	}
	err = rows.Err()
	if err != nil {
		return matchs, err
	}
	return matchs, err
}
