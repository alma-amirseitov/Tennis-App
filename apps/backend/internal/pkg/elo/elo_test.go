package elo

import (
	"math"
	"testing"
)

func almostEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestCalculate_EqualPlayers(t *testing.T) {
	result := Calculate(
		PlayerInfo{Rating: 1200, TotalGames: 25},
		PlayerInfo{Rating: 1200, TotalGames: 25},
	)

	// Equal players → expected 0.50 each, K=32
	// Winner: 1200 + 32*(1.0 - 0.50) = 1216.0
	// Loser:  1200 + 32*(0.0 - 0.50) = 1184.0
	if !almostEqual(result.WinnerNewRating, 1216.0, 0.5) {
		t.Errorf("Winner new rating: expected ~1216.0, got %.1f", result.WinnerNewRating)
	}
	if !almostEqual(result.LoserNewRating, 1184.0, 0.5) {
		t.Errorf("Loser new rating: expected ~1184.0, got %.1f", result.LoserNewRating)
	}
	if result.WinnerDelta <= 0 {
		t.Errorf("Winner delta should be positive, got %.1f", result.WinnerDelta)
	}
	if result.LoserDelta >= 0 {
		t.Errorf("Loser delta should be negative, got %.1f", result.LoserDelta)
	}
}

func TestCalculate_Upset(t *testing.T) {
	// Weak player (1000, 5 games) beats strong player (1400, 50 games)
	result := Calculate(
		PlayerInfo{Rating: 1000, TotalGames: 5},
		PlayerInfo{Rating: 1400, TotalGames: 50},
	)

	// K: 40 (winner), 24 (loser)
	// Expected winner: ~0.091
	// Winner: 1000 + 40*(1.0 - 0.091) ≈ 1036.4
	// Loser:  1400 + 24*(0.0 - 0.909) ≈ 1378.2
	if !almostEqual(result.WinnerNewRating, 1036.4, 1.0) {
		t.Errorf("Winner new rating: expected ~1036.4, got %.1f", result.WinnerNewRating)
	}
	if !almostEqual(result.LoserNewRating, 1378.2, 1.0) {
		t.Errorf("Loser new rating: expected ~1378.2, got %.1f", result.LoserNewRating)
	}

	// Winner delta should be large (upset bonus)
	if result.WinnerDelta < 30 {
		t.Errorf("Upset winner delta should be large (>30), got %.1f", result.WinnerDelta)
	}
}

func TestCalculate_ExpectedResult(t *testing.T) {
	// Strong player (1500, 40 games) beats weak player (1100, 15 games)
	result := Calculate(
		PlayerInfo{Rating: 1500, TotalGames: 40},
		PlayerInfo{Rating: 1100, TotalGames: 15},
	)

	// K: 24 (winner), 32 (loser)
	// Expected winner: ~0.909
	// Winner: 1500 + 24*(1.0 - 0.909) ≈ 1502.2
	// Loser:  1100 + 32*(0.0 - 0.091) ≈ 1097.1
	if !almostEqual(result.WinnerNewRating, 1502.2, 1.0) {
		t.Errorf("Winner new rating: expected ~1502.2, got %.1f", result.WinnerNewRating)
	}
	if !almostEqual(result.LoserNewRating, 1097.1, 1.0) {
		t.Errorf("Loser new rating: expected ~1097.1, got %.1f", result.LoserNewRating)
	}

	// Winner delta should be small (expected result)
	if result.WinnerDelta > 5 {
		t.Errorf("Expected result winner delta should be small (<5), got %.1f", result.WinnerDelta)
	}
}

func TestCalculate_ClampMin(t *testing.T) {
	result := Calculate(
		PlayerInfo{Rating: 200, TotalGames: 50},
		PlayerInfo{Rating: 100, TotalGames: 50},
	)

	if result.LoserNewRating < MinRating {
		t.Errorf("Rating should not go below MinRating (%v), got %.1f", MinRating, result.LoserNewRating)
	}
}

func TestCalculate_ClampMax(t *testing.T) {
	result := Calculate(
		PlayerInfo{Rating: 2990, TotalGames: 50},
		PlayerInfo{Rating: 2000, TotalGames: 50},
	)

	if result.WinnerNewRating > MaxRating {
		t.Errorf("Rating should not go above MaxRating (%v), got %.1f", MaxRating, result.WinnerNewRating)
	}
}

func TestCalculateDoubles(t *testing.T) {
	winners := [2]DoublesPlayerInfo{
		{Rating: 1200, TotalGames: 20},
		{Rating: 1300, TotalGames: 15},
	}
	losers := [2]DoublesPlayerInfo{
		{Rating: 1250, TotalGames: 30},
		{Rating: 1150, TotalGames: 10},
	}

	result := CalculateDoubles(winners, losers)

	// Both winners should gain rating
	if result.Winner1.WinnerDelta <= 0 {
		t.Errorf("Winner1 delta should be positive, got %.1f", result.Winner1.WinnerDelta)
	}
	if result.Winner2.WinnerDelta <= 0 {
		t.Errorf("Winner2 delta should be positive, got %.1f", result.Winner2.WinnerDelta)
	}

	// Both losers should lose rating
	if result.Loser1.LoserDelta >= 0 {
		t.Errorf("Loser1 delta should be negative, got %.1f", result.Loser1.LoserDelta)
	}
	if result.Loser2.LoserDelta >= 0 {
		t.Errorf("Loser2 delta should be negative, got %.1f", result.Loser2.LoserDelta)
	}
}

func TestKFactor(t *testing.T) {
	tests := []struct {
		games    int
		expected int
	}{
		{0, KNew},
		{5, KNew},
		{9, KNew},
		{10, KMedium},
		{20, KMedium},
		{30, KMedium},
		{31, KStable},
		{100, KStable},
	}

	for _, tt := range tests {
		result := kFactor(tt.games)
		if result != tt.expected {
			t.Errorf("kFactor(%d): expected %d, got %d", tt.games, tt.expected, result)
		}
	}
}

func TestGetNTRPLevel(t *testing.T) {
	tests := []struct {
		rating        float64
		expectedLevel string
	}{
		{500, "1.0"},
		{800, "2.0"},
		{1000, "2.5"},
		{1200, "3.0"},
		{1500, "4.0"},
		{2000, "5.5"},
		{2600, "6.5+"},
	}

	for _, tt := range tests {
		level, _ := GetNTRPLevel(tt.rating)
		if level != tt.expectedLevel {
			t.Errorf("GetNTRPLevel(%.0f): expected %s, got %s", tt.rating, tt.expectedLevel, level)
		}
	}
}

func TestClamp(t *testing.T) {
	if clamp(50) != MinRating {
		t.Errorf("clamp(50) should be MinRating, got %.1f", clamp(50))
	}
	if clamp(3500) != MaxRating {
		t.Errorf("clamp(3500) should be MaxRating, got %.1f", clamp(3500))
	}
	if clamp(1500.0) != 1500.0 {
		t.Errorf("clamp(1500.0) should be 1500.0, got %.1f", clamp(1500.0))
	}
}

func TestCalculate_NewPlayer(t *testing.T) {
	// New player (0 games, InitialRating) vs experienced player
	result := Calculate(
		PlayerInfo{Rating: InitialRating, TotalGames: 0},
		PlayerInfo{Rating: 1200, TotalGames: 25},
	)

	// New player K=40, should calibrate faster
	if result.WinnerDelta < 20 {
		t.Errorf("New player winning should get large delta, got %.1f", result.WinnerDelta)
	}
}

func TestCalculate_Symmetry(t *testing.T) {
	// When both players have equal stats, the sum of deltas should be zero
	result := Calculate(
		PlayerInfo{Rating: 1200, TotalGames: 25},
		PlayerInfo{Rating: 1200, TotalGames: 25},
	)

	sum := result.WinnerDelta + result.LoserDelta
	if !almostEqual(sum, 0, 0.5) {
		t.Errorf("Equal player deltas should sum to ~0, got %.1f", sum)
	}
}
