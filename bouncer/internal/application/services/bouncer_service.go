package services

import (
	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/external"
	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/domain/repositories"
	"github.com/MoutonGrognon/robot-game-lua/bouncer/internal/infrastructure/rest"
)

type BouncerService struct {
	botRepo      repositories.BotRepository
	matchRepo    repositories.MatchRepository
	matchmakerMS external.MatchmakerMS
}

func NewBouncerService(botRepo repositories.BotRepository, matchRepo repositories.MatchRepository) BouncerService {
	return BouncerService{
		botRepo:      botRepo,
		matchRepo:    matchRepo,
		matchmakerMS: rest.NewMatchmakerMS(),
	}
}

func (s *BouncerService) AddMatchToQueue(blueName string, redName string) (bool, error) {
	return s.matchmakerMS.AddMatchToQueue(blueName, redName)
}
