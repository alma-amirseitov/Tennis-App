package service

import (
	"context"
	"fmt"
	"time"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// RatingService handles rating and leaderboard business logic
type RatingService struct {
	repo *repository.Queries
}

// NewRatingService creates a new RatingService
func NewRatingService(repo *repository.Queries) *RatingService {
	return &RatingService{repo: repo}
}

// LeaderboardEntry represents a single entry in the leaderboard
type LeaderboardEntry struct {
	Rank      int    `json:"rank"`
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	NTRPLevel *float64 `json:"ntrp_level,omitempty"`
	Rating    float64  `json:"rating"`
	Games     int      `json:"games"`
	Wins      int      `json:"wins"`
	Losses    int      `json:"losses"`
	WinRate   float64  `json:"win_rate"`
}

// RatingHistoryEntry represents a rating change in history
type RatingHistoryEntry struct {
	Date   string  `json:"date"`
	Rating float64 `json:"rating"`
	Change float64 `json:"change"`
	Reason string  `json:"reason,omitempty"`
}

// MyRatingResponse represents the user's rating position
type MyRatingResponse struct {
	Global      GlobalRatingInfo      `json:"global"`
	Communities []CommunityRatingInfo `json:"communities"`
}

// GlobalRatingInfo represents global rating info
type GlobalRatingInfo struct {
	Rating       float64 `json:"rating"`
	Rank         int     `json:"rank"`
	TotalPlayers int64   `json:"total_players"`
}

// CommunityRatingInfo represents community-specific rating info
type CommunityRatingInfo struct {
	CommunityID   string  `json:"community_id"`
	CommunityName string  `json:"community_name"`
	CommunityLogo *string `json:"community_logo,omitempty"`
	Rating        float64 `json:"rating"`
	Games         int     `json:"games"`
	Wins          int     `json:"wins"`
	Losses        int     `json:"losses"`
	Rank          int     `json:"rank"`
}

// ListLeaderboardInput represents input for listing leaderboard
type ListLeaderboardInput struct {
	MinGames int
	Page     int
	PerPage  int
}

// ListRatingHistoryInput represents input for listing rating history
type ListRatingHistoryInput struct {
	CommunityID string
	Period      string // "1m", "3m", "6m", "1y", "all"
	Page        int
	PerPage     int
}

// GetGlobalLeaderboard returns the global leaderboard
func (s *RatingService) GetGlobalLeaderboard(ctx context.Context, input ListLeaderboardInput) ([]LeaderboardEntry, int64, error) {
	page := input.Page
	if page < 1 {
		page = 1
	}
	perPage := input.PerPage
	if perPage < 1 || perPage > 50 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	params := repository.GetGlobalLeaderboardParams{
		ResultOffset: int32(offset),
		ResultLimit:  int32(perPage),
	}
	if input.MinGames > 0 {
		params.MinGames = pgtype.Int4{Int32: int32(input.MinGames), Valid: true}
	}

	countParams := pgtype.Int4{}
	if input.MinGames > 0 {
		countParams = pgtype.Int4{Int32: int32(input.MinGames), Valid: true}
	}

	rows, err := s.repo.GetGlobalLeaderboard(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("get global leaderboard: %w", err)
	}

	total, err := s.repo.CountGlobalLeaderboard(ctx, countParams)
	if err != nil {
		return nil, 0, fmt.Errorf("count global leaderboard: %w", err)
	}

	entries := make([]LeaderboardEntry, 0, len(rows))
	for i, row := range rows {
		entry := LeaderboardEntry{
			Rank:    offset + i + 1,
			UserID:  pgtypeUUIDToStringRequired(row.UserID),
			Rating:  numericToFloat(row.GlobalRating),
			Games:   int(row.TotalGames),
			Wins:    int(row.TotalWins),
			Losses:  int(row.TotalLosses),
			WinRate: numericToFloat(row.WinRate),
		}
		if row.FirstName.Valid {
			entry.FirstName = row.FirstName.String
		}
		if row.LastName.Valid {
			entry.LastName = row.LastName.String
		}
		if row.AvatarUrl.Valid {
			entry.AvatarURL = &row.AvatarUrl.String
		}
		if row.NtrpLevel.Valid {
			ntrp := numericToFloat(row.NtrpLevel)
			entry.NTRPLevel = &ntrp
		}
		entries = append(entries, entry)
	}

	return entries, total, nil
}

// GetCommunityLeaderboard returns the community leaderboard
func (s *RatingService) GetCommunityLeaderboard(ctx context.Context, communityID uuid.UUID, input ListLeaderboardInput) ([]LeaderboardEntry, error) {
	page := input.Page
	if page < 1 {
		page = 1
	}
	perPage := input.PerPage
	if perPage < 1 || perPage > 50 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	minGames := input.MinGames
	if minGames < 0 {
		minGames = 0
	}

	rows, err := s.repo.GetCommunityLeaderboard(ctx, repository.GetCommunityLeaderboardParams{
		CommunityID:         uuidToPgtype(communityID),
		CommunityGamesCount: pgtype.Int4{Int32: int32(minGames), Valid: true},
		Limit:               int32(perPage),
		Offset:              int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("get community leaderboard: %w", err)
	}

	entries := make([]LeaderboardEntry, 0, len(rows))
	for i, row := range rows {
		entry := LeaderboardEntry{
			Rank:    offset + i + 1,
			UserID:  pgtypeUUIDToStringRequired(row.UserID),
			Rating:  numericToFloat(row.CommunityRating),
			Games:   int(row.CommunityGamesCount.Int32),
			Wins:    int(row.CommunityWins.Int32),
			Losses:  int(row.CommunityLosses.Int32),
			WinRate: float64(row.WinRate),
		}
		if row.FirstName.Valid {
			entry.FirstName = row.FirstName.String
		}
		if row.LastName.Valid {
			entry.LastName = row.LastName.String
		}
		if row.AvatarUrl.Valid {
			entry.AvatarURL = &row.AvatarUrl.String
		}
		if row.NtrpLevel.Valid {
			ntrp := numericToFloat(row.NtrpLevel)
			entry.NTRPLevel = &ntrp
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// GetMyRatingHistory returns the user's rating history for graphing
func (s *RatingService) GetMyRatingHistory(ctx context.Context, userID uuid.UUID, input ListRatingHistoryInput) ([]RatingHistoryEntry, error) {
	page := input.Page
	if page < 1 {
		page = 1
	}
	perPage := input.PerPage
	if perPage < 1 || perPage > 100 {
		perPage = 50
	}
	offset := (page - 1) * perPage

	params := repository.GetRatingHistoryParams{
		UserID:       uuidToPgtype(userID),
		ResultOffset: int32(offset),
		ResultLimit:  int32(perPage),
	}

	if input.CommunityID != "" {
		commUUID, err := uuid.Parse(input.CommunityID)
		if err == nil {
			params.CommunityID = uuidToPgtype(commUUID)
		}
	}

	// Calculate "since" based on period
	if input.Period != "" && input.Period != "all" {
		since := periodToTime(input.Period)
		if since != nil {
			params.Since = pgtype.Timestamptz{Time: *since, Valid: true}
		}
	}

	rows, err := s.repo.GetRatingHistory(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("get rating history: %w", err)
	}

	entries := make([]RatingHistoryEntry, 0, len(rows))
	for _, row := range rows {
		entry := RatingHistoryEntry{
			Date:   row.CreatedAt.Time.Format(time.RFC3339),
			Rating: numericToFloat(row.RatingAfter),
			Change: numericToFloat(row.Change),
		}
		if row.Reason.Valid {
			entry.Reason = row.Reason.String
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// GetMyRating returns the user's rating position globally and in communities
func (s *RatingService) GetMyRating(ctx context.Context, userID uuid.UUID) (*MyRatingResponse, error) {
	// Get global position
	position, err := s.repo.GetUserRatingPosition(ctx, uuidToPgtype(userID))
	if err != nil {
		return nil, fmt.Errorf("get user rating position: %w", err)
	}

	// Get community ratings
	communityRatings, err := s.repo.GetUserCommunityRatings(ctx, uuidToPgtype(userID))
	if err != nil {
		return nil, fmt.Errorf("get user community ratings: %w", err)
	}

	resp := &MyRatingResponse{
		Global: GlobalRatingInfo{
			Rating:       numericToFloat(position.GlobalRating),
			Rank:         int(position.Rank),
			TotalPlayers: position.TotalPlayers,
		},
		Communities: make([]CommunityRatingInfo, 0, len(communityRatings)),
	}

	for _, cr := range communityRatings {
		info := CommunityRatingInfo{
			CommunityID:   pgtypeUUIDToStringRequired(cr.CommunityID),
			CommunityName: cr.CommunityName,
			Rating:        numericToFloat(cr.CommunityRating),
			Games:         int(cr.CommunityGamesCount.Int32),
			Wins:          int(cr.CommunityWins.Int32),
			Losses:        int(cr.CommunityLosses.Int32),
			Rank:          int(cr.Rank),
		}
		if cr.CommunityLogo.Valid {
			info.CommunityLogo = &cr.CommunityLogo.String
		}
		resp.Communities = append(resp.Communities, info)
	}

	return resp, nil
}

// GetMyStats returns the user's global statistics
func (s *RatingService) GetMyStats(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	stats, err := s.repo.GetUserStats(ctx, uuidToPgtype(userID))
	if err != nil {
		// Return default stats if no stats found
		return map[string]interface{}{
			"total_games":        0,
			"total_wins":         0,
			"total_losses":       0,
			"win_rate":           0.0,
			"singles_games":      0,
			"singles_wins":       0,
			"doubles_games":      0,
			"doubles_wins":       0,
			"current_streak":     0,
			"best_streak":        0,
			"tournaments_played": 0,
		}, nil
	}

	return map[string]interface{}{
		"total_games":        stats.TotalGames.Int32,
		"total_wins":         stats.TotalWins.Int32,
		"total_losses":       stats.TotalLosses.Int32,
		"win_rate":           numericToFloat(stats.WinRate),
		"singles_games":      stats.SinglesGames.Int32,
		"singles_wins":       stats.SinglesWins.Int32,
		"doubles_games":      stats.DoublesGames.Int32,
		"doubles_wins":       stats.DoublesWins.Int32,
		"current_streak":     stats.CurrentStreak.Int32,
		"best_streak":        stats.BestStreak.Int32,
		"tournaments_played": stats.TournamentsPlayed.Int32,
	}, nil
}

// periodToTime converts a period string to a time.Time pointer
func periodToTime(period string) *time.Time {
	now := time.Now()
	var t time.Time
	switch period {
	case "1m":
		t = now.AddDate(0, -1, 0)
	case "3m":
		t = now.AddDate(0, -3, 0)
	case "6m":
		t = now.AddDate(0, -6, 0)
	case "1y":
		t = now.AddDate(-1, 0, 0)
	default:
		return nil
	}
	return &t
}
