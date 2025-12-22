package services

import (
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/entities"
	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/external"
	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/repositories"
	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/infrastructure/rest"
)

type BouncerService struct {
	botRepo          repositories.BotRepository
	matchRepo        repositories.MatchRepository
	rankingRepo      repositories.RankingRepository
	matchmakerMS     external.MatchmakerMS
	highlightedMatch *entities.Match
}

func NewBouncerService(botRepo repositories.BotRepository, matchRepo repositories.MatchRepository, rankingRepo repositories.RankingRepository) *BouncerService {
	return &BouncerService{
		botRepo:      botRepo,
		matchRepo:    matchRepo,
		rankingRepo:  rankingRepo,
		matchmakerMS: rest.NewMatchmakerMS(),
	}
}

func (s *BouncerService) AddMatchToQueue(blueName string, redName string) (bool, error) {
	return s.matchmakerMS.AddMatchToQueue(blueName, redName)
}

func (s *BouncerService) GetMatch(matchId uuid.UUID) (entities.Match, error) {
	return s.matchRepo.GetById(matchId)
}

func (s *BouncerService) GetSummaries(start int, size int) ([]entities.MatchSummary, int, int, int, error) {
	return s.matchRepo.GetSummaries(start, size)
}

func (s *BouncerService) GetRanking() ([]entities.Rank, error) {
	return s.rankingRepo.GetRanking()
}

func (s *BouncerService) GetHighlightedMatch() (*entities.Match, error) {
	defer s.GetHighlightedMatchDebounced()
	if s.highlightedMatch != nil {
		return s.highlightedMatch, nil
	}
	return s.GetHighlightedMatchDebounced()
}

func (s *BouncerService) GetHighlightedMatchDebounced() (*entities.Match, error) {
	highlightDuration := time.Hour * 12
	if s.highlightedMatch == nil || s.highlightedMatch.Date.Add(highlightDuration).Before(time.Now()) {
		summaries, err := s.matchRepo.GetRecentSummaries(time.Now().Add(-highlightDuration))
		if err != nil {
			return nil, err
		}
		ranks, err := s.rankingRepo.GetRanking()
		if err != nil {
			return nil, err
		}
		if len(summaries) == 0 {
			return nil, nil
		}
		matchId := summaries[0].Id
		bestHighlightScore := math.MinInt
		eloById := map[uuid.UUID]int{}
		for _, summary := range summaries {
			blueElo, blueEloFound := eloById[summary.BlueBotId]
			redElo, redEloFound := eloById[summary.RedBotId]
			if !blueEloFound {
				for _, rank := range ranks {
					if rank.BotId == summary.BlueBotId {
						eloById[rank.BotId] = blueElo
						blueElo, blueEloFound = eloById[summary.BlueBotId]
						break
					}
				}
			}
			if !redEloFound {
				for _, rank := range ranks {
					if rank.BotId == summary.RedBotId {
						eloById[rank.BotId] = redElo
						redElo, redEloFound = eloById[summary.RedBotId]
						break
					}
				}
			}
			if !blueEloFound || !redEloFound {
				continue
			}
			lowElo := blueElo
			highElo := redElo
			if blueElo > redElo {
				lowElo = redElo
				highElo = blueElo
			}
			scoreDifference := summary.BlueScore - summary.RedScore
			// Arbitrary formula to get high elo but low elo difference,
			// while having a close match (low score difference)
			highlightScore := 3*lowElo - highElo - scoreDifference*scoreDifference
			if highlightScore > bestHighlightScore {
				bestHighlightScore = highlightScore
				matchId = summary.Id
			}
		}
		match, err := s.GetMatch(matchId)
		if err != nil {
			return nil, err
		}
		return &match, nil
	}
	return s.highlightedMatch, nil
}
