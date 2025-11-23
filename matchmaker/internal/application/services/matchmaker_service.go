package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/parquet-go/parquet-go/compress/zstd"

	"github.com/MoutonGrognon/robot-game-lua/matchmaker/internal/domain/entities"
	"github.com/MoutonGrognon/robot-game-lua/matchmaker/internal/domain/external"
	"github.com/MoutonGrognon/robot-game-lua/matchmaker/internal/domain/repositories"
	"github.com/MoutonGrognon/robot-game-lua/matchmaker/internal/infrastructure/rest"
	"github.com/MoutonGrognon/robot-game-lua/rgcore/rgconst"
	"github.com/MoutonGrognon/robot-game-lua/rgcore/rgentities"
	"github.com/MoutonGrognon/robot-game-lua/rgcore/rgutils"
)

const MATCH_TIMEOUT = 2 *
	// Convert from milliseconds to nanoseconds
	1000000 *
	// Init time
	((rgconst.BOT_INIT_TIME_BUDGET *
		// number of bots per wave
		2 * rgconst.SPAWN_COUNT *
		// number of waves
		rgconst.MAX_TURN / rgconst.SPAWN_DELAY) +

		// Action time
		(rgconst.BOT_ACTION_TIME_BUDGET *
			// duration of a wave
			rgconst.SPAWN_DELAY *
			// sum of the max number of bots per wave
			((rgconst.MAX_TURN / rgconst.SPAWN_DELAY) *
				((rgconst.MAX_TURN / rgconst.SPAWN_DELAY) + 1) / 2) *
			2 * rgconst.SPAWN_COUNT))

type MatchmakerService struct {
	botRepo           repositories.BotRepository
	matchRepo         repositories.MatchRepository
	refereeMS         external.RefereeMS
	matchQueue        entities.MatchQueue
	isRunning         bool
	currentMatch      *entities.PendingMatch
	debounceTimer     *time.Timer
	matchMu           sync.Mutex
	forcedRankedMatch bool
	rankedMatchTimer  *time.Timer
	rankedMu          sync.Mutex
}

func NewMatchmakerService(botRepo repositories.BotRepository, matchRepo repositories.MatchRepository) MatchmakerService {
	matchmakerService := MatchmakerService{
		botRepo:           botRepo,
		matchRepo:         matchRepo,
		refereeMS:         rest.NewRefereeMS(),
		matchQueue:        entities.NewMatchQueue(),
		isRunning:         false,
		forcedRankedMatch: true,
		currentMatch:      &entities.PendingMatch{},
	}
	matchmakerService.forceDebouncedRankedMatch()
	go func() {
		err := matchmakerService.StartDebouncedMatch()
		for err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Println("Wait 1s and retry ...")
			time.Sleep(1 * time.Second)
			err = matchmakerService.StartDebouncedMatch()
		}
		fmt.Println("Auto-matching successfully started")
	}()
	return matchmakerService
}

func (s *MatchmakerService) forceDebouncedRankedMatch() {
	s.rankedMu.Lock()
	s.forcedRankedMatch = true
	if s.rankedMatchTimer != nil {
		s.rankedMatchTimer.Stop()
	}
	// Matchs are very unlikely to reach a duration of MATCH_TIMEOUT
	// so it should leave quite some time for unranked matchs
	s.rankedMatchTimer = time.AfterFunc(MATCH_TIMEOUT, func() {
		s.forceDebouncedRankedMatch()
	})
	s.rankedMu.Unlock()
}

func (s *MatchmakerService) printGrid(currentGameState map[int]rgentities.BotState) {
	gameStateAsStr := ""
	for i := 0; i < rgconst.GRID_SIZE; i++ {
		for j := 0; j < rgconst.GRID_SIZE; j++ {
			tile := " "
			if rgutils.GetLocationType(j, i) == rgconst.OBSTACLE {
				tile = "#"
			}
			gameStateAsStr += tile + " "
		}
		gameStateAsStr += "\n"
	}
	blueCount := 0
	redCount := 0
	for _, botState := range currentGameState {
		tile := "O"
		blueCount += 1
		if botState.Bot.PlayerId == rgconst.RED_ID {
			tile = "X"
			blueCount -= 1
			redCount += 1
		}
		tileIndex := ((2*rgconst.GRID_SIZE+1)*botState.Bot.Y + (2 * botState.Bot.X))
		gameStateAsStr = gameStateAsStr[:tileIndex] + tile + gameStateAsStr[tileIndex+1:]
	}
	gameStateAsStr = strings.ReplaceAll(gameStateAsStr, "# ", "\033[40m  \033[47m")
	gameStateAsStr = strings.ReplaceAll(gameStateAsStr, "X ", "\033[41m  \033[47m")
	gameStateAsStr = strings.ReplaceAll(gameStateAsStr, "O ", "\033[46m  \033[47m")
	gameStateAsStr = strings.ReplaceAll(gameStateAsStr, "\n", "\033[0m\n\033[47m")
	fmt.Printf("%d - %d\n", blueCount, redCount)
	fmt.Printf("\033[47m%s\033[0m\n", gameStateAsStr)
}

func (s *MatchmakerService) SaveMatch(matchId uuid.UUID, game []map[int]rgentities.BotState) error {
	for i, state := range game {
		fmt.Printf("turn %d\n", i+1)
		s.printGrid(state)
	}
	score1 := 0
	score2 := 0
	for _, botState := range game[len(game)-1] {
		if botState.Bot.PlayerId == rgconst.BLUE_ID {
			score1 += 1
		} else {
			score2 += 1
		}
	}
	fmt.Printf("%v - %v\n", score1, score2)
	s.matchMu.Lock()
	defer s.matchMu.Unlock()
	if matchId != s.currentMatch.Id {
		err := errors.New("Corrupted Match ID")
		fmt.Printf("Error: %v\n", err)
		return err
	}
	s.isRunning = false
	var err error
	defer func() {
		go s.StartDebouncedMatch()
	}()
	jsonGame, err := json.Marshal(game)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	codec := zstd.Codec{}
	compressedGame, err := codec.Encode(nil, jsonGame)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	err = s.matchRepo.Save(entities.Match{
		Id:             matchId,
		BotId1:         s.currentMatch.BotId1,
		BotId2:         s.currentMatch.BotId2,
		BotName1:       s.currentMatch.BotName1,
		BotName2:       s.currentMatch.BotName2,
		UserName1:      s.currentMatch.UserName1,
		UserName2:      s.currentMatch.UserName2,
		Date:           time.Now(),
		CompressedGame: compressedGame,
		Score1:         score1,
		Score2:         score2,
		Ranked:         s.currentMatch.Ranked,
	})
	// TODO: update bots elo if match was ranked
	return err
}

func (s *MatchmakerService) CancelMatch(matchId uuid.UUID, err error) error {
	// TODO: cancel match
	return nil
}

func (s *MatchmakerService) KillMatch() error {
	s.matchMu.Lock()
	defer func() {
		s.isRunning = false
		s.matchMu.Unlock()
	}()
	return s.refereeMS.KillMatch()
}

func (s *MatchmakerService) CreateMatchFromNames(blueName string, redName string, ranked bool) (entities.PendingMatch, error) {
	pendingMatch := entities.PendingMatch{}
	blueId, err := s.botRepo.GetIdFromName(blueName)
	if err != nil {
		return pendingMatch, err
	}
	redId, err := s.botRepo.GetIdFromName(redName)
	if err != nil {
		return pendingMatch, err
	}
	blueUserName, err := s.botRepo.GetUserNameFromBotId(blueId)
	if err != nil {
		return pendingMatch, err
	}
	redUserName, err := s.botRepo.GetUserNameFromBotId(redId)
	if err != nil {
		return pendingMatch, err
	}
	pendingMatch = entities.PendingMatch{
		Id:        uuid.New(),
		BotId1:    blueId,
		BotId2:    redId,
		BotName1:  blueName,
		BotName2:  redName,
		UserName1: blueUserName,
		UserName2: redUserName,
		Ranked:    ranked,
	}
	return pendingMatch, nil
}

func (s *MatchmakerService) AddMatchToQueue(blueName string, redName string) (bool, error) {
	// TODO: better system to handle queue size and check on elements added
	if s.matchQueue.IsFull() {
		return false, nil
	}
	pendingMatch, err := s.CreateMatchFromNames(blueName, redName, false)
	if err != nil {
		return false, err
	}
	added := s.matchQueue.Push(pendingMatch)
	go s.StartDebouncedMatch()
	return added, nil
}

func (s *MatchmakerService) StartMatch(pendingMatch entities.PendingMatch) error {
	*s.currentMatch = pendingMatch
	s.isRunning = true
	err := s.refereeMS.StartMatch(pendingMatch.Id, pendingMatch.BotId1, pendingMatch.BotId2)
	if err != nil {
		s.isRunning = false
	}
	return err
}

func (s *MatchmakerService) StartDebouncedMatch() error {
	var err error
	s.matchMu.Lock()
	defer s.matchMu.Unlock()
	if s.isRunning {
		return nil
	}
	if s.debounceTimer != nil {
		s.debounceTimer.Stop()
	}
	var pendingMatch entities.PendingMatch
	if s.forcedRankedMatch {
		s.rankedMu.Lock()
		s.forcedRankedMatch = false
		s.rankedMu.Unlock()
		// TODO: WIP choose bots based on their elo
		blueName := "random"
		redName := "random"
		pendingMatch, err = s.CreateMatchFromNames(blueName, redName, true)
	} else {
		// force next match to be ranked
		s.forceDebouncedRankedMatch()
		if s.matchQueue.IsEmpty() {
			fmt.Printf("Sleep for at most %vs\n", MATCH_TIMEOUT/1000000000)
			time.AfterFunc(MATCH_TIMEOUT, func() {
				s.StartDebouncedMatch()
			})
		}
		pendingMatch, err = s.matchQueue.Pop()
	}
	if err != nil {
		return err
	}
	s.debounceTimer = time.AfterFunc(MATCH_TIMEOUT, func() {
		s.KillMatch()
		fmt.Println("Kill match because it took too long")
		err := s.StartDebouncedMatch()
		for err != nil {
			fmt.Printf("Error: %v\n", err)
			if s.matchQueue.IsEmpty() {
				return
			}
			fmt.Println("Try next match in queue")
			err = s.StartDebouncedMatch()
		}
	})
	return s.StartMatch(pendingMatch)
}
