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
	stmt, err := mr.db.Prepare("SELECT * FROM matches WHERE id=$1")
	if err != nil {
		return match, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&match.Id, &match.BlueBotId, &match.RedBotId, &match.BlueBotName, &match.RedBotName, &match.BlueUserName, &match.RedUserName, &match.Date, &match.CompressedGame, &match.BlueScore, &match.RedScore, &match.Ranked)
	if err != nil {
		return match, err
	}
	return match, err
}

func (mr *MatchRepository) GetSummaries(start int, size int) ([]entities.MatchSummary, int, int, int, error) {
	var matches []entities.MatchSummary
	total := 0
	stmt, err := mr.db.Prepare("SELECT id, blueBotId, redBotId, blueBotName, redBotName, blueUserName, redUserName, date, blueScore, redScore, ranked FROM (" +
		"SELECT ROW_NUMBER() OVER (ORDER BY date DESC) as rowNum, * FROM matches" +
		") WHERE (rowNum>=$1 AND rowNum<$2) ORDER BY rowNum")
	if err != nil {
		return matches, start, size, total, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(start, start+size)
	if err != nil {
		return matches, start, size, total, err
	}
	defer rows.Close()
	for rows.Next() {
		var match entities.MatchSummary
		err = rows.Scan(&match.Id, &match.BlueBotId, &match.RedBotId, &match.BlueBotName, &match.RedBotName, &match.BlueUserName, &match.RedUserName, &match.Date, &match.BlueScore, &match.RedScore, &match.Ranked)
		if err != nil {
			return matches, start, size, total, err
		}
		matches = append(matches, match)
	}
	err = rows.Err()
	if err != nil {
		return matches, start, size, total, err
	}
	row := mr.db.QueryRow("SELECT count FROM matches_count")
	row.Scan(&total)
	return matches, start, size, total, err
}

func (mr *MatchRepository) GetRecentSummaries(dateThreshold time.Time) ([]entities.MatchSummary, error) {
	var matches []entities.MatchSummary
	stmt, err := mr.db.Prepare("SELECT id, blueBotId, redBotId, blueBotName, redBotName, blueUserName, redUserName, date, blueScore, redScore, ranked FROM matches WHERE date>=$1")
	if err != nil {
		return matches, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(dateThreshold)
	if err != nil {
		return matches, err
	}
	defer rows.Close()
	for rows.Next() {
		var match entities.MatchSummary
		err = rows.Scan(&match.Id, &match.BlueBotId, &match.RedBotId, &match.BlueBotName, &match.RedBotName, &match.BlueUserName, &match.RedUserName, &match.Date, &match.BlueScore, &match.RedScore, &match.Ranked)
		if err != nil {
			return matches, err
		}
		matches = append(matches, match)
	}
	err = rows.Err()
	if err != nil {
		return matches, err
	}
	return matches, err
}

func (mr *MatchRepository) DeleteOldMatches(threshold int) (int, error) {
	stmt, err := mr.db.Prepare("SELECT date FROM (" +
		"SELECT ROW_NUMBER() OVER (ORDER BY date DESC) as rowNum, * FROM matches" +
		") WHERE rowNum=$1 ORDER BY rowNum")
	if err != nil {
		return 0, err
	}
	var dateThreshold time.Time
	err = stmt.QueryRow(threshold).Scan(&dateThreshold)
	stmt.Close()
	if err != nil {
		return 0, err
	}
	stmt, err = mr.db.Prepare("DELETE FROM matches WHERE date<$1")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(dateThreshold)
	if err != nil {
		return 0, err
	}
	deletedRowsCount, err := res.RowsAffected()
	if err != nil {
		return int(deletedRowsCount), err
	}
	return int(deletedRowsCount), nil
}
