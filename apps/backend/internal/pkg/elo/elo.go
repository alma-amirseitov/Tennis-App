package elo

import (
	"math"
)

const (
	InitialRating = 1000.0
	MinRating     = 100.0
	MaxRating     = 3000.0
	KNew          = 40  // < 10 games (fast calibration)
	KMedium       = 32  // 10-30 games
	KStable       = 24  // > 30 games (stable)
	ScaleFactor   = 400.0
)

// PlayerInfo contains the data needed for ELO calculation
type PlayerInfo struct {
	Rating     float64
	TotalGames int
}

// RatingChange represents the calculated rating changes after a match
type RatingChange struct {
	WinnerNewRating float64
	LoserNewRating  float64
	WinnerDelta     float64 // always positive
	LoserDelta      float64 // always negative
}

// Calculate computes new ratings after a singles match
func Calculate(winner, loser PlayerInfo) RatingChange {
	// Expected scores
	expWinner := expectedScore(winner.Rating, loser.Rating)
	expLoser := 1.0 - expWinner

	// K-factors (per player)
	kWinner := kFactor(winner.TotalGames)
	kLoser := kFactor(loser.TotalGames)

	// New ratings
	winnerDelta := float64(kWinner) * (1.0 - expWinner)
	loserDelta := float64(kLoser) * (0.0 - expLoser)

	newWinner := clamp(winner.Rating + winnerDelta)
	newLoser := clamp(loser.Rating + loserDelta)

	return RatingChange{
		WinnerNewRating: newWinner,
		LoserNewRating:  newLoser,
		WinnerDelta:     newWinner - winner.Rating,
		LoserDelta:      newLoser - loser.Rating,
	}
}

// DoublesPlayerInfo contains info for a doubles team member
type DoublesPlayerInfo struct {
	Rating     float64
	TotalGames int
}

// DoublesRatingChange contains rating changes for all 4 players in doubles
type DoublesRatingChange struct {
	Winner1 RatingChange // team1 player 1 (winner side)
	Winner2 RatingChange // team1 player 2 (winner side)
	Loser1  RatingChange // team2 player 1 (loser side)
	Loser2  RatingChange // team2 player 2 (loser side)
}

// CalculateDoubles computes ratings for a doubles match
// Each player is rated individually against the average rating of the opposing team
func CalculateDoubles(winners [2]DoublesPlayerInfo, losers [2]DoublesPlayerInfo) DoublesRatingChange {
	avgWinnerRating := (winners[0].Rating + winners[1].Rating) / 2
	avgLoserRating := (losers[0].Rating + losers[1].Rating) / 2

	// Winner 1 vs avg loser team
	w1 := Calculate(
		PlayerInfo{Rating: winners[0].Rating, TotalGames: winners[0].TotalGames},
		PlayerInfo{Rating: avgLoserRating, TotalGames: 0},
	)

	// Winner 2 vs avg loser team
	w2 := Calculate(
		PlayerInfo{Rating: winners[1].Rating, TotalGames: winners[1].TotalGames},
		PlayerInfo{Rating: avgLoserRating, TotalGames: 0},
	)

	// Loser 1 vs avg winner team
	l1 := Calculate(
		PlayerInfo{Rating: avgWinnerRating, TotalGames: 0},
		PlayerInfo{Rating: losers[0].Rating, TotalGames: losers[0].TotalGames},
	)

	// Loser 2 vs avg winner team
	l2 := Calculate(
		PlayerInfo{Rating: avgWinnerRating, TotalGames: 0},
		PlayerInfo{Rating: losers[1].Rating, TotalGames: losers[1].TotalGames},
	)

	return DoublesRatingChange{
		Winner1: w1,
		Winner2: w2,
		Loser1:  l1,
		Loser2:  l2,
	}
}

// NTRPLevel represents a mapping from rating to NTRP skill level
type NTRPLevel struct {
	Min   float64
	Max   float64
	Level string
	Name  string
}

// NTRPLevels maps rating ranges to NTRP skill levels
var NTRPLevels = []NTRPLevel{
	{100, 599, "1.0", "Начинающий"},
	{600, 799, "1.5", "Начинающий+"},
	{800, 999, "2.0", "Новичок"},
	{1000, 1149, "2.5", "Новичок+"},
	{1150, 1299, "3.0", "Любитель"},
	{1300, 1449, "3.5", "Любитель+"},
	{1450, 1599, "4.0", "Продвинутый"},
	{1600, 1799, "4.5", "Продвинутый+"},
	{1800, 1999, "5.0", "Опытный"},
	{2000, 2299, "5.5", "Опытный+"},
	{2300, 2599, "6.0", "Полупрофессионал"},
	{2600, 3000, "6.5+", "Профессионал"},
}

// GetNTRPLevel returns the NTRP level string and label for a given rating
func GetNTRPLevel(rating float64) (string, string) {
	for _, l := range NTRPLevels {
		if rating >= l.Min && rating <= l.Max {
			return l.Level, l.Name
		}
	}
	return "1.0", "Начинающий"
}

func expectedScore(ratingA, ratingB float64) float64 {
	return 1.0 / (1.0 + math.Pow(10, (ratingB-ratingA)/ScaleFactor))
}

func kFactor(totalGames int) int {
	switch {
	case totalGames < 10:
		return KNew
	case totalGames <= 30:
		return KMedium
	default:
		return KStable
	}
}

func clamp(rating float64) float64 {
	if rating < MinRating {
		return MinRating
	}
	if rating > MaxRating {
		return MaxRating
	}
	return math.Round(rating*10) / 10 // round to 1 decimal
}
