# ELO Rating Algorithm
## Система рейтинга теннисной платформы

---

## 1. Overview

Modified ELO system для любительского тенниса. Рейтинг определяет позицию игрока в глобальном и community-лидербордах.

## 2. Constants

```go
const (
    InitialRating     = 1000.0  // Стартовый рейтинг нового игрока
    MinRating         = 100.0   // Минимальный возможный рейтинг
    MaxRating         = 3000.0  // Максимальный возможный рейтинг
    KFactorNew        = 40      // < 10 игр (быстрая калибровка)
    KFactorMedium     = 32      // 10-30 игр
    KFactorStable     = 24      // > 30 игр (стабильный)
    ScaleFactor       = 400.0   // Стандартный ELO scale factor
)
```

## 3. Algorithm

### Step 1: Expected Score
```
E_winner = 1 / (1 + 10^((R_loser - R_winner) / 400))
E_loser  = 1 - E_winner
```

### Step 2: K-Factor Selection
```
K = 40  if total_games < 10    (новичок, быстрая калибровка)
K = 32  if total_games 10-30   (средний)
K = 24  if total_games > 30    (стабильный)
```

K-фактор определяется для КАЖДОГО игрока ОТДЕЛЬНО по его количеству игр.

### Step 3: New Rating
```
R_new_winner = R_old_winner + K_winner * (1.0 - E_winner)
R_new_loser  = R_old_loser  + K_loser  * (0.0 - E_loser)
```

### Step 4: Clamp
```
R_new = clamp(R_new, MinRating, MaxRating)
```

## 4. Implementation

```go
package elo

import (
    "math"
)

const (
    InitialRating = 1000.0
    MinRating     = 100.0
    MaxRating     = 3000.0
    KNew          = 40
    KMedium       = 32
    KStable       = 24
    ScaleFactor   = 400.0
)

// Result represents a match outcome
type Result struct {
    WinnerOldRating float64
    LoserOldRating  float64
    WinnerGames     int  // total games played by winner
    LoserGames      int  // total games played by loser
}

// RatingChange represents the calculated rating changes
type RatingChange struct {
    WinnerNewRating float64
    LoserNewRating  float64
    WinnerDelta     float64 // always positive
    LoserDelta      float64 // always negative
}

// Calculate computes new ratings after a singles match
func Calculate(r Result) RatingChange {
    // Expected scores
    expWinner := expectedScore(r.WinnerOldRating, r.LoserOldRating)
    expLoser := 1.0 - expWinner

    // K-factors (per player)
    kWinner := kFactor(r.WinnerGames)
    kLoser := kFactor(r.LoserGames)

    // New ratings
    winnerDelta := float64(kWinner) * (1.0 - expWinner)
    loserDelta := float64(kLoser) * (0.0 - expLoser)

    newWinner := clamp(r.WinnerOldRating + winnerDelta)
    newLoser := clamp(r.LoserOldRating + loserDelta)

    return RatingChange{
        WinnerNewRating: newWinner,
        LoserNewRating:  newLoser,
        WinnerDelta:     newWinner - r.WinnerOldRating,
        LoserDelta:      newLoser - r.LoserOldRating,
    }
}

// CalculateDoubles computes ratings for doubles match
// Each player gets individual rating change based on average opponent rating
func CalculateDoubles(team1, team2 [2]Result) [4]RatingChange {
    // Average ratings
    avgTeam1 := (team1[0].WinnerOldRating + team1[1].WinnerOldRating) / 2
    avgTeam2 := (team2[0].LoserOldRating + team2[1].LoserOldRating) / 2

    var changes [4]RatingChange
    // Each winner vs average loser team rating
    for i := 0; i < 2; i++ {
        c := Calculate(Result{
            WinnerOldRating: team1[i].WinnerOldRating,
            LoserOldRating:  avgTeam2,
            WinnerGames:     team1[i].WinnerGames,
            LoserGames:      0, // not used for winner calc
        })
        changes[i] = c
    }
    // Each loser vs average winner team rating
    for i := 0; i < 2; i++ {
        c := Calculate(Result{
            WinnerOldRating: avgTeam1,
            LoserOldRating:  team2[i].LoserOldRating,
            WinnerGames:     0,
            LoserGames:      team2[i].LoserGames,
        })
        changes[i+2] = c
    }
    return changes
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
```

## 5. Examples

### Example 1: Equal players
```
Winner: 1200 (25 games) vs Loser: 1200 (25 games)
Expected: 0.50 vs 0.50
K: 32 vs 32
Winner: 1200 + 32*(1.0 - 0.50) = 1216.0 (+16.0)
Loser:  1200 + 32*(0.0 - 0.50) = 1184.0 (-16.0)
```

### Example 2: Upset (weak beats strong)
```
Winner: 1000 (5 games) vs Loser: 1400 (50 games)
Expected: 0.091 vs 0.909
K: 40 vs 24
Winner: 1000 + 40*(1.0 - 0.091) = 1036.4 (+36.4)
Loser:  1400 + 24*(0.0 - 0.909) = 1378.2 (-21.8)
```

### Example 3: Expected result (strong beats weak)
```
Winner: 1500 (40 games) vs Loser: 1100 (15 games)
Expected: 0.909 vs 0.091
K: 24 vs 32
Winner: 1500 + 24*(1.0 - 0.909) = 1502.2 (+2.2)
Loser:  1100 + 32*(0.0 - 0.091) = 1097.1 (-2.9)
```

## 6. NTRP Level Mapping

Rating → display level for user profile:

```go
var NTRPLevels = []struct {
    Min   float64
    Max   float64
    Level string
    Name  string
}{
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

func GetNTRPLevel(rating float64) (string, string) {
    for _, l := range NTRPLevels {
        if rating >= l.Min && rating <= l.Max {
            return l.Level, l.Name
        }
    }
    return "1.0", "Начинающий"
}
```

## 7. When Rating is Calculated

1. Both players submit match result
2. Results MATCH (or second player confirms first player's result)
3. System calls `elo.Calculate()` with both players' current ratings and game counts
4. New ratings saved to `users.rating_score` and `rating_history` table
5. Community leaderboards updated
6. Push notification sent to both players with delta

## 8. Edge Cases

| Case | Handling |
|------|---------|
| Player has 0 games | Use InitialRating (1000), K=40 |
| Self-match (same player) | Reject at service layer |
| Result dispute (different scores) | No rating change until resolved |
| Result not confirmed in 48h | Auto-confirm based on submitter's data |
| Rating below 100 | Clamp to MinRating |
| Rating above 3000 | Clamp to MaxRating |
| Doubles match | Use average opponent team rating |
| Tournament match | Same formula, no special handling |
| Walkover / forfeit | Winner gets minimum change (+2), loser gets normal loss |

## 9. Database

```sql
-- Rating change is stored in rating_history
INSERT INTO rating_history (user_id, match_id, community_id, old_rating, new_rating, delta)
VALUES ($1, $2, $3, $4, $5, $6);

-- User rating is updated
UPDATE users SET rating_score = $1, games_played = games_played + 1 WHERE id = $2;

-- Community leaderboard is a VIEW
-- See docs/03-DATABASE-SCHEMA.sql for community_leaderboard view
```
