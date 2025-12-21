package services

import (
<<<<<<< HEAD
	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/entities"
	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/external"
	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/repositories"
	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/infrastructure/rest"
=======
	"time"

	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/entities"
	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/external"
	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/repositories"
	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/infrastructure/rest"
>>>>>>> d97a808 (feat: WIP add frontend + add highlighted match endpoint)
	"github.com/google/uuid"
)

type BouncerService struct {
	botRepo          repositories.BotRepository
	matchRepo        repositories.MatchRepository
	matchmakerMS     external.MatchmakerMS
	highlightedMatch *entities.Match
}

func NewBouncerService(botRepo repositories.BotRepository, matchRepo repositories.MatchRepository) *BouncerService {
	return &BouncerService{
		botRepo:      botRepo,
		matchRepo:    matchRepo,
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

func (s *BouncerService) GetHighlightedMatch() (*entities.Match, error) {
	defer s.GetHighlightedMatchDebounced()
	if s.highlightedMatch != nil {
		return s.highlightedMatch, nil
	}
	return s.GetHighlightedMatchDebounced()
}

func (s *BouncerService) GetHighlightedMatchDebounced() (*entities.Match, error) {
	if s.highlightedMatch == nil || s.highlightedMatch.Date.Add(time.Hour*12).Before(time.Now()) {
		// TODO: find some metrics to discover interesting matchs
		summaries, _, _, _, err := s.GetSummaries(1, 1)
		if err != nil {
			return nil, err
		}
		if len(summaries) == 0 {
			return nil, nil
		}
		matchId := summaries[0].Id
		match, err := s.GetMatch(matchId)
		if err != nil {
			return nil, err
		}
		return &match, nil
	}
	return s.highlightedMatch, nil
}
