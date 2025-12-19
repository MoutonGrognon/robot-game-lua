package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
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

const (
	// Arbitrary value increasing the elo variation for bots
	// that have played less games than the chosen threshold
	CONVERGENCE_THRESHOLD = 20
	K                     = 20.0
	// Increased elo varation factor for new bots
	K_BOOSTED = K * 8
	// TODO: get from conf
	DEFAULT_ELO = 1000

	MATCH_TIMEOUT = 2 *
		// Convert from milliseconds to nanoseconds
		1000 * 1000 *
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

	// TODO: get from conf
	// TODO: rework the system to allow more control on
	// ranked/casual matches frequency, repartition, etc
	RANKED_MATCH_INTERVALL = 10 * (1000 * 1000 * 1000) // 10s
)

type MatchmakerService struct {
	botRepo           repositories.BotRepository
	matchRepo         repositories.MatchRepository
	rankingRepo       repositories.RankingRepository
	refereeMS         external.RefereeMS
	matchQueue        *entities.MatchQueue
	isRunning         bool
	currentMatch      *entities.PendingMatch
	debounceTimer     *time.Timer
	matchMu           sync.Mutex
	forcedRankedMatch bool
	rankedMatchTimer  *time.Timer
	rankedMu          sync.Mutex
	ranks             []entities.Rank
}

func NewMatchmakerService(botRepo repositories.BotRepository, matchRepo repositories.MatchRepository, rankingRepo repositories.RankingRepository) *MatchmakerService {
	matchmakerService := &MatchmakerService{
		botRepo:           botRepo,
		matchRepo:         matchRepo,
		rankingRepo:       rankingRepo,
		refereeMS:         rest.NewRefereeMS(),
		matchQueue:        entities.NewMatchQueue(),
		isRunning:         false,
		forcedRankedMatch: true,
		currentMatch:      &entities.PendingMatch{},
	}
	ranks, err := matchmakerService.rankingRepo.GetRanking()
	for err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Wait 1s and retry ...")
		time.Sleep(1 * time.Second)
		ranks, err = matchmakerService.rankingRepo.GetRanking()
	}
	matchmakerService.ranks = ranks
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
	s.rankedMatchTimer = time.AfterFunc(RANKED_MATCH_INTERVALL, func() {
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

func (s *MatchmakerService) MatchEsperance(elo1 int, elo2 int) float64 {
	return 1.0 / (1.0 + math.Pow(10.0, float64(elo2-elo1)/400.0))
}

// Factor K to modulate elo variation.
// New bots have a high K to converge faster to their "true" elo
// Old bots have a lower K to have a more stable elo
func (s *MatchmakerService) dynamicK(matchCount int) float64 {
	if matchCount >= CONVERGENCE_THRESHOLD {
		return K
	}
	return (float64(matchCount)*K + float64(CONVERGENCE_THRESHOLD-matchCount)*K_BOOSTED) / CONVERGENCE_THRESHOLD
}

// Factor to reduce elo variation against new bots,
// to avoid excessive win/loss based on an elo not
// representative of the true level of the opponent.
// It also greatly slows elo loss at a very low elo.
func (s *MatchmakerService) fairnessBalancingFactor(elo int, matchCount int, opponentMatchCount int) float64 {
	if elo < K {
		return 0.0
	}
	f := 1.0
	if elo < 400 {
		f *= float64(elo) / 400
	}
	if opponentMatchCount >= CONVERGENCE_THRESHOLD || matchCount <= opponentMatchCount {
		return f
	}
	if matchCount > CONVERGENCE_THRESHOLD {
		f *= (float64(opponentMatchCount) / CONVERGENCE_THRESHOLD)
	}
	f *= (float64(opponentMatchCount) / CONVERGENCE_THRESHOLD)
	return f

}

func (s *MatchmakerService) matchCount(rank entities.Rank) int {
	return rank.WinCount + rank.LossCount + rank.DrawCount
}

func (s *MatchmakerService) chooseBots() (string, string, error) {
	if len(s.ranks) < 2 {
		return "", "", errors.New("Not enough bots loaded for a match")
	}
	blueIndex := 0
	redIndex := 1

	// TODO: Change the system because as soon as it is possible
	// to add a new bot after the initial bots have been there
	// for a little while, the new bot will be in all matches
	// until it reaches the same match count as the other,
	// which will be quite boring

	// Choose the bot that played the least amount of matchs
	for rankIndex, rank := range s.ranks {
		if s.matchCount(rank) < s.matchCount(s.ranks[blueIndex]) {
			blueIndex = rankIndex
		}
	}
	// Choose the most relevant opponent based on
	// - elo variation (high elo variation is better)
	// - match count (low match count is better)
	bestRelevanceScore := 0.0
	blueElo := s.ranks[blueIndex].Elo
	blueMatchCount := s.matchCount(s.ranks[blueIndex])
	for rankIndex, rank := range s.ranks {
		if rankIndex == blueIndex {
			continue
		}
		candidateElo := rank.Elo
		candidateMatchCount := s.matchCount(rank)
		blueK := s.dynamicK(blueMatchCount) * s.fairnessBalancingFactor(blueElo, blueMatchCount, candidateMatchCount)
		redK := s.dynamicK(candidateMatchCount) * s.fairnessBalancingFactor(candidateElo, candidateMatchCount, blueMatchCount)

		// Can be used for the probability of a win* and the expected
		// normalized variation of elo (the elo system is built to
		// compensate a high win probabilty with a low reward so for
		// probability of win* P (with a probability of loss* 1-P)
		// the normalized gain is set K*(1-P) and
		// the normalized loss is set to K*(-P)
		// such that the expected variation is
		// E = K*P*(1-P)+K*(1-P)(-P) = 0
		// *assuming the probability of a draw is 0. By counting a draw
		// as both half a win and half a loss the math checks out and
		// and there is no need to consider draws separatly
		blueWinEsperance := s.MatchEsperance(blueElo, candidateElo)
		eloVariationEsperance := (blueK + redK) * blueWinEsperance * (1 - blueWinEsperance)
		lowMatchCountFactor := CONVERGENCE_THRESHOLD
		if blueMatchCount < candidateMatchCount {
			lowMatchCountFactor += blueMatchCount
		} else {
			lowMatchCountFactor += candidateMatchCount
		}
		highMatchCount := 10 * CONVERGENCE_THRESHOLD
		if lowMatchCountFactor > highMatchCount {
			lowMatchCountFactor = highMatchCount
		}
		// Arbitrary formula that benefits to match causing high
		// elo variation (probably improving convergence speed ?)
		// while giving less chances to bots that have already played
		// a lot of matchs
		relevanceScore := eloVariationEsperance / float64(lowMatchCountFactor)
		if relevanceScore > bestRelevanceScore {
			redIndex = rankIndex
			bestRelevanceScore = relevanceScore
		}
	}
	return s.ranks[blueIndex].BotName, s.ranks[redIndex].BotName, nil
}

func (s *MatchmakerService) UpdateRanking(match entities.Match) error {
	blueRankIndex := -1
	redRankIndex := -1
	for i, rank := range s.ranks {
		if rank.BotId == match.BotId1 {
			blueRankIndex = i
		}
		if rank.BotId == match.BotId2 {
			redRankIndex = i
		}
	}
	if blueRankIndex < 0 {
		blueRankIndex = len(s.ranks)
		s.ranks = append(s.ranks, entities.Rank{
			Id:        uuid.New(),
			BotId:     match.BotId1,
			BotName:   match.BotName1,
			Elo:       DEFAULT_ELO,
			WinCount:  0,
			DrawCount: 0,
			LossCount: 0,
		})
	}
	if redRankIndex < 0 {
		redRankIndex = len(s.ranks)
		s.ranks = append(s.ranks, entities.Rank{
			Id:        uuid.New(),
			BotId:     match.BotId2,
			BotName:   match.BotName2,
			Elo:       DEFAULT_ELO,
			WinCount:  0,
			DrawCount: 0,
			LossCount: 0,
		})
	}
	var res float64
	blueMatchCount := s.matchCount(s.ranks[blueRankIndex])
	redMatchCount := s.matchCount(s.ranks[redRankIndex])
	blueElo := s.ranks[blueRankIndex].Elo
	redElo := s.ranks[redRankIndex].Elo
	blueK := s.dynamicK(blueMatchCount) * s.fairnessBalancingFactor(blueElo, blueMatchCount, redMatchCount)
	redK := s.dynamicK(redMatchCount) * s.fairnessBalancingFactor(redElo, redMatchCount, blueMatchCount)
	if match.Score1 > match.Score2 {
		res = 1.0
		s.ranks[blueRankIndex].WinCount += 1
		s.ranks[redRankIndex].LossCount += 1
	} else if match.Score1 < match.Score2 {
		res = 0.0
		s.ranks[blueRankIndex].LossCount += 1
		s.ranks[redRankIndex].WinCount += 1
	} else {
		res = 0.5
		s.ranks[blueRankIndex].DrawCount += 1
		s.ranks[redRankIndex].DrawCount += 1
	}
	s.ranks[blueRankIndex].Elo += int(math.Round(blueK * (res - s.MatchEsperance(s.ranks[blueRankIndex].Elo, s.ranks[redRankIndex].Elo))))
	s.ranks[redRankIndex].Elo -= int(math.Round(redK * (res - s.MatchEsperance(s.ranks[blueRankIndex].Elo, s.ranks[redRankIndex].Elo))))
	var err error
	err = s.rankingRepo.UpdateRank(s.ranks[blueRankIndex])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	err = s.rankingRepo.UpdateRank(s.ranks[redRankIndex])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	return nil
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
	// TODO: WIP begin transaction
	match := entities.Match{
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
	}
	err = s.matchRepo.Save(match)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	if s.currentMatch.Ranked {
		err = s.UpdateRanking(match)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
	}
	// TODO: WIP end transaction
	return nil
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
		blueName, redName, noChosenBotsError := s.chooseBots()
		if noChosenBotsError != nil {
			return noChosenBotsError
		}
		pendingMatch, err = s.CreateMatchFromNames(blueName, redName, true)
	} else {
		// force next match to be ranked
		s.forceDebouncedRankedMatch()
		if s.matchQueue.IsEmpty() {
			fmt.Printf("Sleep for at most %vs\n", RANKED_MATCH_INTERVALL/1000000000)
			time.AfterFunc(RANKED_MATCH_INTERVALL, func() {
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
