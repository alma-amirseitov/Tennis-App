package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/pkg/elo"
	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MatchService handles match business logic
type MatchService struct {
	repo          *repository.Queries
	pool          *pgxpool.Pool
	notifications *NotificationService
}

// NewMatchService creates a new MatchService
func NewMatchService(repo *repository.Queries, pool *pgxpool.Pool, notifications *NotificationService) *MatchService {
	return &MatchService{
		repo:          repo,
		pool:          pool,
		notifications: notifications,
	}
}

// SetScore represents a single set's score
type SetScore struct {
	Player1  int  `json:"p1"`
	Player2  int  `json:"p2"`
	Tiebreak *int `json:"tiebreak,omitempty"`
}

// SubmitResultInput represents input for submitting a match result
type SubmitResultInput struct {
	WinnerID string     `json:"winner_id"`
	Score    []SetScore `json:"score"`
}

// ConfirmResultInput represents input for confirming/disputing a match result
type ConfirmResultInput struct {
	Action string `json:"action"` // "confirm" or "dispute"
	Reason string `json:"reason,omitempty"`
}

// AdminConfirmInput represents input for admin confirming a disputed match
type AdminConfirmInput struct {
	WinnerID string     `json:"winner_id"`
	Score    []SetScore `json:"score"`
	Note     string     `json:"note,omitempty"`
}

// RatingChangeInfo represents a player's rating change in a match response
type RatingChangeInfo struct {
	Before float64 `json:"before"`
	After  float64 `json:"after"`
	Change float64 `json:"change"`
}

// MatchResponse represents a match in API responses
type MatchResponse struct {
	ID             string  `json:"id"`
	EventID        *string `json:"event_id,omitempty"`
	CommunityID    *string `json:"community_id,omitempty"`
	Player1ID      string  `json:"player1_id"`
	Player2ID      string  `json:"player2_id"`
	Player1Partner *string `json:"player1_partner_id,omitempty"`
	Player2Partner *string `json:"player2_partner_id,omitempty"`
	Composition    string  `json:"composition"`
	Score          any     `json:"score,omitempty"`
	WinnerID       *string `json:"winner_id,omitempty"`
	ResultStatus   string  `json:"result_status"`
	SubmittedBy    *string `json:"submitted_by,omitempty"`
	ConfirmedBy    *string `json:"confirmed_by,omitempty"`
	SubmittedAt    *string `json:"submitted_at,omitempty"`
	ConfirmedAt    *string `json:"confirmed_at,omitempty"`
	DisputeReason  *string `json:"dispute_reason,omitempty"`
	RatingChanges  *struct {
		Player1 *RatingChangeInfo `json:"player1,omitempty"`
		Player2 *RatingChangeInfo `json:"player2,omitempty"`
	} `json:"rating_changes,omitempty"`
	RoundName     *string `json:"round_name,omitempty"`
	RoundNumber   *int    `json:"round_number,omitempty"`
	CourtNumber   *int    `json:"court_number,omitempty"`
	ScheduledTime *string `json:"scheduled_time,omitempty"`
	PlayedAt      *string `json:"played_at,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

// ListMyMatchesInput represents input for listing matches
type ListMyMatchesInput struct {
	CommunityID string
	OpponentID  string
	Result      string // "win", "loss", "all"
	Page        int
	PerPage     int
}

// GetByID returns a match by ID
func (s *MatchService) GetByID(ctx context.Context, matchID uuid.UUID) (*MatchResponse, error) {
	match, err := s.repo.GetMatchByID(ctx, uuidToPgtype(matchID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrMatchNotFound
		}
		return nil, fmt.Errorf("get match by id: %w", err)
	}
	resp := buildMatchResponse(match)
	return &resp, nil
}

// SubmitResult submits a match result (first player action)
func (s *MatchService) SubmitResult(ctx context.Context, userID uuid.UUID, matchID uuid.UUID, input SubmitResultInput) (*MatchResponse, error) {
	// Get the match
	match, err := s.repo.GetMatchByID(ctx, uuidToPgtype(matchID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrMatchNotFound
		}
		return nil, fmt.Errorf("get match: %w", err)
	}

	// Check if user is a player in this match
	if !isPlayerInMatch(match, userID) {
		return nil, ErrForbidden.WithMessage("You are not a player in this match")
	}

	// Check if result is already submitted
	if match.ResultStatus.Valid && match.ResultStatus.ResultStatus != repository.ResultStatusPending {
		// If there's already a winner set, result was already submitted
		if match.WinnerID.Valid {
			return nil, ErrResultAlreadySubmit
		}
	}

	// Validate winner is a player
	winnerUUID, err := uuid.Parse(input.WinnerID)
	if err != nil {
		return nil, ErrValidation.WithMessage("Invalid winner_id")
	}
	if !isPlayerInMatch(match, winnerUUID) {
		return nil, ErrValidation.WithMessage("Winner must be a player in this match")
	}

	// Validate score
	if len(input.Score) == 0 {
		return nil, ErrValidation.WithMessage("Score is required")
	}

	scoreJSON, err := json.Marshal(input.Score)
	if err != nil {
		return nil, fmt.Errorf("marshal score: %w", err)
	}

	updated, err := s.repo.SubmitMatchResult(ctx, repository.SubmitMatchResultParams{
		Score:       scoreJSON,
		WinnerID:    uuidToPgtype(winnerUUID),
		SubmittedBy: uuidToPgtype(userID),
		ID:          uuidToPgtype(matchID),
	})
	if err != nil {
		return nil, fmt.Errorf("submit match result: %w", err)
	}

	// Notify opponent(s) that result is pending confirmation
	if s.notifications != nil {
		s.notifyMatchResultPending(ctx, match, userID, matchID)
	}

	resp := buildMatchResponse(updated)
	return &resp, nil
}

// ConfirmResult confirms or disputes a match result (second player action)
func (s *MatchService) ConfirmResult(ctx context.Context, userID uuid.UUID, matchID uuid.UUID, input ConfirmResultInput) (*MatchResponse, error) {
	// Get the match
	match, err := s.repo.GetMatchByID(ctx, uuidToPgtype(matchID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrMatchNotFound
		}
		return nil, fmt.Errorf("get match: %w", err)
	}

	// Check if user is a player in this match
	if !isPlayerInMatch(match, userID) {
		return nil, ErrForbidden.WithMessage("You are not a player in this match")
	}

	// Result must be pending
	if !match.ResultStatus.Valid || match.ResultStatus.ResultStatus != repository.ResultStatusPending {
		return nil, ErrValidation.WithMessage("Match result is not pending confirmation")
	}

	// Confirmer must be the other player (not the submitter)
	submitterID, _ := uuid.FromBytes(match.SubmittedBy.Bytes[:])
	if userID == submitterID {
		return nil, ErrForbidden.WithMessage("You cannot confirm your own result submission")
	}

	switch input.Action {
	case "confirm":
		return s.confirmMatch(ctx, userID, match)
	case "dispute":
		return s.disputeMatch(ctx, matchID, input.Reason)
	default:
		return nil, ErrValidation.WithMessage("Action must be 'confirm' or 'dispute'")
	}
}

// confirmMatch handles the match confirmation and ELO calculation in a transaction
func (s *MatchService) confirmMatch(ctx context.Context, confirmerID uuid.UUID, match repository.Match) (*MatchResponse, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	// Get winner and loser
	winnerID, _ := uuid.FromBytes(match.WinnerID.Bytes[:])
	player1ID, _ := uuid.FromBytes(match.Player1ID.Bytes[:])
	player2ID, _ := uuid.FromBytes(match.Player2ID.Bytes[:])

	var loserID uuid.UUID
	if winnerID == player1ID {
		loserID = player2ID
	} else {
		loserID = player1ID
	}

	// Get player ratings
	winnerData, err := qtx.GetUserForRating(ctx, uuidToPgtype(winnerID))
	if err != nil {
		return nil, fmt.Errorf("get winner rating: %w", err)
	}
	loserData, err := qtx.GetUserForRating(ctx, uuidToPgtype(loserID))
	if err != nil {
		return nil, fmt.Errorf("get loser rating: %w", err)
	}

	winnerRating := numericToFloat(winnerData.GlobalRating)
	loserRating := numericToFloat(loserData.GlobalRating)

	// Get total games for K-factor
	winnerGames, _ := qtx.GetPlayerTotalGames(ctx, uuidToPgtype(winnerID))
	loserGames, _ := qtx.GetPlayerTotalGames(ctx, uuidToPgtype(loserID))

	// Calculate ELO
	ratingChange := elo.Calculate(
		elo.PlayerInfo{Rating: winnerRating, TotalGames: int(winnerGames)},
		elo.PlayerInfo{Rating: loserRating, TotalGames: int(loserGames)},
	)

	// Determine which player is P1 and P2 for rating storage
	var p1RatingBefore, p1RatingAfter, p2RatingBefore, p2RatingAfter float64
	if winnerID == player1ID {
		p1RatingBefore = winnerRating
		p1RatingAfter = ratingChange.WinnerNewRating
		p2RatingBefore = loserRating
		p2RatingAfter = ratingChange.LoserNewRating
	} else {
		p1RatingBefore = loserRating
		p1RatingAfter = ratingChange.LoserNewRating
		p2RatingBefore = winnerRating
		p2RatingAfter = ratingChange.WinnerNewRating
	}

	// 1. Update match with ratings
	confirmed, err := qtx.ConfirmMatch(ctx, repository.ConfirmMatchParams{
		ConfirmedBy:         uuidToPgtype(confirmerID),
		Player1RatingBefore: floatToNumeric(p1RatingBefore),
		Player1RatingAfter:  floatToNumeric(p1RatingAfter),
		Player2RatingBefore: floatToNumeric(p2RatingBefore),
		Player2RatingAfter:  floatToNumeric(p2RatingAfter),
		ID:                  match.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("confirm match: %w", err)
	}

	// 2. Update winner rating
	if err := qtx.UpdateUserRating(ctx, repository.UpdateUserRatingParams{
		NewRating: floatToNumeric(ratingChange.WinnerNewRating),
		UserID:    uuidToPgtype(winnerID),
	}); err != nil {
		return nil, fmt.Errorf("update winner rating: %w", err)
	}

	// 3. Update loser rating
	if err := qtx.UpdateUserRating(ctx, repository.UpdateUserRatingParams{
		NewRating: floatToNumeric(ratingChange.LoserNewRating),
		UserID:    uuidToPgtype(loserID),
	}); err != nil {
		return nil, fmt.Errorf("update loser rating: %w", err)
	}

	// 4. Update NTRP levels
	winnerNTRP, _ := elo.GetNTRPLevel(ratingChange.WinnerNewRating)
	loserNTRP, _ := elo.GetNTRPLevel(ratingChange.LoserNewRating)
	s.updateNTRPIfChanged(ctx, qtx, winnerID, winnerNTRP)
	s.updateNTRPIfChanged(ctx, qtx, loserID, loserNTRP)

	isSingles := match.Composition == repository.PlayerCompositionSingles

	// 5. Update player stats
	if err := qtx.UpsertPlayerStatsGlobal(ctx, repository.UpsertPlayerStatsGlobalParams{
		UserID:    uuidToPgtype(winnerID),
		IsWinner:  true,
		IsSingles: isSingles,
	}); err != nil {
		return nil, fmt.Errorf("upsert winner stats: %w", err)
	}
	if err := qtx.UpsertPlayerStatsGlobal(ctx, repository.UpsertPlayerStatsGlobalParams{
		UserID:    uuidToPgtype(loserID),
		IsWinner:  false,
		IsSingles: isSingles,
	}); err != nil {
		return nil, fmt.Errorf("upsert loser stats: %w", err)
	}

	// 6. Insert rating history (global)
	if _, err := qtx.InsertRatingHistory(ctx, repository.InsertRatingHistoryParams{
		UserID:       uuidToPgtype(winnerID),
		RatingBefore: floatToNumeric(winnerRating),
		RatingAfter:  floatToNumeric(ratingChange.WinnerNewRating),
		Change:       floatToNumeric(ratingChange.WinnerDelta),
		MatchID:      match.ID,
		Reason:       pgtype.Text{String: "match_result", Valid: true},
	}); err != nil {
		return nil, fmt.Errorf("insert winner rating history: %w", err)
	}
	if _, err := qtx.InsertRatingHistory(ctx, repository.InsertRatingHistoryParams{
		UserID:       uuidToPgtype(loserID),
		RatingBefore: floatToNumeric(loserRating),
		RatingAfter:  floatToNumeric(ratingChange.LoserNewRating),
		Change:       floatToNumeric(ratingChange.LoserDelta),
		MatchID:      match.ID,
		Reason:       pgtype.Text{String: "match_result", Valid: true},
	}); err != nil {
		return nil, fmt.Errorf("insert loser rating history: %w", err)
	}

	// 7. Update community member stats if match is within a community
	if match.CommunityID.Valid {
		if err := qtx.UpdateCommunityMemberStats(ctx, repository.UpdateCommunityMemberStatsParams{
			NewRating:   floatToNumeric(ratingChange.WinnerNewRating),
			IsWinner:    true,
			CommunityID: match.CommunityID,
			UserID:      uuidToPgtype(winnerID),
		}); err != nil {
			slog.Warn("failed to update community winner stats", "error", err)
		}
		if err := qtx.UpdateCommunityMemberStats(ctx, repository.UpdateCommunityMemberStatsParams{
			NewRating:   floatToNumeric(ratingChange.LoserNewRating),
			IsWinner:    false,
			CommunityID: match.CommunityID,
			UserID:      uuidToPgtype(loserID),
		}); err != nil {
			slog.Warn("failed to update community loser stats", "error", err)
		}

		// Insert community-specific rating history
		if _, err := qtx.InsertRatingHistory(ctx, repository.InsertRatingHistoryParams{
			UserID:       uuidToPgtype(winnerID),
			CommunityID:  match.CommunityID,
			RatingBefore: floatToNumeric(winnerRating),
			RatingAfter:  floatToNumeric(ratingChange.WinnerNewRating),
			Change:       floatToNumeric(ratingChange.WinnerDelta),
			MatchID:      match.ID,
			Reason:       pgtype.Text{String: "match_result", Valid: true},
		}); err != nil {
			slog.Warn("failed to insert community winner rating history", "error", err)
		}
		if _, err := qtx.InsertRatingHistory(ctx, repository.InsertRatingHistoryParams{
			UserID:       uuidToPgtype(loserID),
			CommunityID:  match.CommunityID,
			RatingBefore: floatToNumeric(loserRating),
			RatingAfter:  floatToNumeric(ratingChange.LoserNewRating),
			Change:       floatToNumeric(ratingChange.LoserDelta),
			MatchID:      match.ID,
			Reason:       pgtype.Text{String: "match_result", Valid: true},
		}); err != nil {
			slog.Warn("failed to insert community loser rating history", "error", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	resp := buildMatchResponse(confirmed)

	// Send rating_changed notifications (best-effort)
	if s.notifications != nil {
		s.notifyRatingChanged(ctx, winnerID, loserID, match.ID, ratingChange)
	}

	return &resp, nil
}

// disputeMatch marks a match result as disputed
func (s *MatchService) disputeMatch(ctx context.Context, matchID uuid.UUID, reason string) (*MatchResponse, error) {
	disputed, err := s.repo.DisputeMatch(ctx, repository.DisputeMatchParams{
		DisputeReason: pgtype.Text{String: reason, Valid: reason != ""},
		ID:            uuidToPgtype(matchID),
	})
	if err != nil {
		return nil, fmt.Errorf("dispute match: %w", err)
	}
	resp := buildMatchResponse(disputed)
	return &resp, nil
}

// AdminConfirm handles admin confirmation of a disputed match
func (s *MatchService) AdminConfirm(ctx context.Context, adminID uuid.UUID, matchID uuid.UUID, input AdminConfirmInput) (*MatchResponse, error) {
	// Get the match
	match, err := s.repo.GetMatchByID(ctx, uuidToPgtype(matchID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrMatchNotFound
		}
		return nil, fmt.Errorf("get match: %w", err)
	}

	// Must be disputed
	if !match.ResultStatus.Valid || match.ResultStatus.ResultStatus != repository.ResultStatusDisputed {
		return nil, ErrValidation.WithMessage("Match must be in disputed status for admin confirmation")
	}

	// Validate winner
	winnerUUID, err := uuid.Parse(input.WinnerID)
	if err != nil {
		return nil, ErrValidation.WithMessage("Invalid winner_id")
	}
	if !isPlayerInMatch(match, winnerUUID) {
		return nil, ErrValidation.WithMessage("Winner must be a player in this match")
	}

	scoreJSON, err := json.Marshal(input.Score)
	if err != nil {
		return nil, fmt.Errorf("marshal score: %w", err)
	}

	// Run the same ELO flow as confirm, but in admin context
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	player1ID, _ := uuid.FromBytes(match.Player1ID.Bytes[:])
	player2ID, _ := uuid.FromBytes(match.Player2ID.Bytes[:])

	var loserID uuid.UUID
	if winnerUUID == player1ID {
		loserID = player2ID
	} else {
		loserID = player1ID
	}

	// Get ratings
	winnerData, err := qtx.GetUserForRating(ctx, uuidToPgtype(winnerUUID))
	if err != nil {
		return nil, fmt.Errorf("get winner rating: %w", err)
	}
	loserData, err := qtx.GetUserForRating(ctx, uuidToPgtype(loserID))
	if err != nil {
		return nil, fmt.Errorf("get loser rating: %w", err)
	}

	winnerRating := numericToFloat(winnerData.GlobalRating)
	loserRating := numericToFloat(loserData.GlobalRating)

	winnerGames, _ := qtx.GetPlayerTotalGames(ctx, uuidToPgtype(winnerUUID))
	loserGames, _ := qtx.GetPlayerTotalGames(ctx, uuidToPgtype(loserID))

	ratingChange := elo.Calculate(
		elo.PlayerInfo{Rating: winnerRating, TotalGames: int(winnerGames)},
		elo.PlayerInfo{Rating: loserRating, TotalGames: int(loserGames)},
	)

	var p1RatingBefore, p1RatingAfter, p2RatingBefore, p2RatingAfter float64
	if winnerUUID == player1ID {
		p1RatingBefore = winnerRating
		p1RatingAfter = ratingChange.WinnerNewRating
		p2RatingBefore = loserRating
		p2RatingAfter = ratingChange.LoserNewRating
	} else {
		p1RatingBefore = loserRating
		p1RatingAfter = ratingChange.LoserNewRating
		p2RatingBefore = winnerRating
		p2RatingAfter = ratingChange.WinnerNewRating
	}

	confirmed, err := qtx.AdminConfirmMatch(ctx, repository.AdminConfirmMatchParams{
		Score:               scoreJSON,
		WinnerID:            uuidToPgtype(winnerUUID),
		ConfirmedBy:         uuidToPgtype(adminID),
		Player1RatingBefore: floatToNumeric(p1RatingBefore),
		Player1RatingAfter:  floatToNumeric(p1RatingAfter),
		Player2RatingBefore: floatToNumeric(p2RatingBefore),
		Player2RatingAfter:  floatToNumeric(p2RatingAfter),
		ID:                  uuidToPgtype(matchID),
	})
	if err != nil {
		return nil, fmt.Errorf("admin confirm match: %w", err)
	}

	// Update ratings, stats, history (same as confirmMatch)
	if err := qtx.UpdateUserRating(ctx, repository.UpdateUserRatingParams{
		NewRating: floatToNumeric(ratingChange.WinnerNewRating),
		UserID:    uuidToPgtype(winnerUUID),
	}); err != nil {
		return nil, fmt.Errorf("update winner rating: %w", err)
	}
	if err := qtx.UpdateUserRating(ctx, repository.UpdateUserRatingParams{
		NewRating: floatToNumeric(ratingChange.LoserNewRating),
		UserID:    uuidToPgtype(loserID),
	}); err != nil {
		return nil, fmt.Errorf("update loser rating: %w", err)
	}

	winnerNTRP, _ := elo.GetNTRPLevel(ratingChange.WinnerNewRating)
	loserNTRP, _ := elo.GetNTRPLevel(ratingChange.LoserNewRating)
	s.updateNTRPIfChanged(ctx, qtx, winnerUUID, winnerNTRP)
	s.updateNTRPIfChanged(ctx, qtx, loserID, loserNTRP)

	isSingles := match.Composition == repository.PlayerCompositionSingles

	if err := qtx.UpsertPlayerStatsGlobal(ctx, repository.UpsertPlayerStatsGlobalParams{
		UserID:    uuidToPgtype(winnerUUID),
		IsWinner:  true,
		IsSingles: isSingles,
	}); err != nil {
		return nil, fmt.Errorf("upsert winner stats: %w", err)
	}
	if err := qtx.UpsertPlayerStatsGlobal(ctx, repository.UpsertPlayerStatsGlobalParams{
		UserID:    uuidToPgtype(loserID),
		IsWinner:  false,
		IsSingles: isSingles,
	}); err != nil {
		return nil, fmt.Errorf("upsert loser stats: %w", err)
	}

	if _, err := qtx.InsertRatingHistory(ctx, repository.InsertRatingHistoryParams{
		UserID:       uuidToPgtype(winnerUUID),
		RatingBefore: floatToNumeric(winnerRating),
		RatingAfter:  floatToNumeric(ratingChange.WinnerNewRating),
		Change:       floatToNumeric(ratingChange.WinnerDelta),
		MatchID:      match.ID,
		Reason:       pgtype.Text{String: "admin_confirmed", Valid: true},
	}); err != nil {
		return nil, fmt.Errorf("insert winner rating history: %w", err)
	}
	if _, err := qtx.InsertRatingHistory(ctx, repository.InsertRatingHistoryParams{
		UserID:       uuidToPgtype(loserID),
		RatingBefore: floatToNumeric(loserRating),
		RatingAfter:  floatToNumeric(ratingChange.LoserNewRating),
		Change:       floatToNumeric(ratingChange.LoserDelta),
		MatchID:      match.ID,
		Reason:       pgtype.Text{String: "admin_confirmed", Valid: true},
	}); err != nil {
		return nil, fmt.Errorf("insert loser rating history: %w", err)
	}

	if match.CommunityID.Valid {
		if err := qtx.UpdateCommunityMemberStats(ctx, repository.UpdateCommunityMemberStatsParams{
			NewRating:   floatToNumeric(ratingChange.WinnerNewRating),
			IsWinner:    true,
			CommunityID: match.CommunityID,
			UserID:      uuidToPgtype(winnerUUID),
		}); err != nil {
			slog.Warn("failed to update community winner stats", "error", err)
		}
		if err := qtx.UpdateCommunityMemberStats(ctx, repository.UpdateCommunityMemberStatsParams{
			NewRating:   floatToNumeric(ratingChange.LoserNewRating),
			IsWinner:    false,
			CommunityID: match.CommunityID,
			UserID:      uuidToPgtype(loserID),
		}); err != nil {
			slog.Warn("failed to update community loser stats", "error", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	resp := buildMatchResponse(confirmed)

	// Send rating_changed notifications for admin-confirmed result as well
	if s.notifications != nil {
		s.notifyRatingChanged(ctx, winnerUUID, loserID, match.ID, ratingChange)
	}

	return &resp, nil
}

// ListMyMatches returns the user's match history
func (s *MatchService) ListMyMatches(ctx context.Context, userID uuid.UUID, input ListMyMatchesInput) ([]MatchResponse, int64, error) {
	page := input.Page
	if page < 1 {
		page = 1
	}
	perPage := input.PerPage
	if perPage < 1 || perPage > 50 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	params := repository.ListMyMatchesParams{
		UserID:       uuidToPgtype(userID),
		ResultOffset: int32(offset),
		ResultLimit:  int32(perPage),
	}
	countParams := repository.CountMyMatchesParams{
		UserID: uuidToPgtype(userID),
	}

	if input.CommunityID != "" {
		commUUID, err := uuid.Parse(input.CommunityID)
		if err == nil {
			params.CommunityID = uuidToPgtype(commUUID)
			countParams.CommunityID = uuidToPgtype(commUUID)
		}
	}
	if input.OpponentID != "" {
		oppUUID, err := uuid.Parse(input.OpponentID)
		if err == nil {
			params.OpponentID = uuidToPgtype(oppUUID)
			countParams.OpponentID = uuidToPgtype(oppUUID)
		}
	}
	if input.Result != "" && input.Result != "all" {
		params.ResultFilter = pgtype.Text{String: input.Result, Valid: true}
		countParams.ResultFilter = pgtype.Text{String: input.Result, Valid: true}
	}

	matches, err := s.repo.ListMyMatches(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("list my matches: %w", err)
	}

	total, err := s.repo.CountMyMatches(ctx, countParams)
	if err != nil {
		return nil, 0, fmt.Errorf("count my matches: %w", err)
	}

	results := make([]MatchResponse, 0, len(matches))
	for _, m := range matches {
		results = append(results, buildMatchResponse(m))
	}

	return results, total, nil
}

// updateNTRPIfChanged updates NTRP level if it has changed
func (s *MatchService) updateNTRPIfChanged(ctx context.Context, qtx *repository.Queries, userID uuid.UUID, newNTRP string) {
	// Parse NTRP string to float for storage
	var ntrpFloat float64
	fmt.Sscanf(newNTRP, "%f", &ntrpFloat)
	if ntrpFloat > 0 {
		if err := qtx.UpdateUserNTRPLevel(ctx, repository.UpdateUserNTRPLevelParams{
			NtrpLevel: floatToNumeric(ntrpFloat),
			UserID:    uuidToPgtype(userID),
		}); err != nil {
			slog.Warn("failed to update NTRP level", "user_id", userID, "error", err)
		}
	}
}

// Helper functions

func isPlayerInMatch(match repository.Match, userID uuid.UUID) bool {
	pgID := uuidToPgtype(userID)
	if match.Player1ID == pgID || match.Player2ID == pgID {
		return true
	}
	if match.Player1PartnerID.Valid && match.Player1PartnerID == pgID {
		return true
	}
	if match.Player2PartnerID.Valid && match.Player2PartnerID == pgID {
		return true
	}
	return false
}

func uuidToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func floatToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	n.Scan(fmt.Sprintf("%.2f", f))
	return n
}

func pgtypeUUIDToString(u pgtype.UUID) *string {
	if !u.Valid {
		return nil
	}
	id, err := uuid.FromBytes(u.Bytes[:])
	if err != nil {
		return nil
	}
	s := id.String()
	return &s
}

func pgtypeUUIDToStringRequired(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	id, err := uuid.FromBytes(u.Bytes[:])
	if err != nil {
		return ""
	}
	return id.String()
}

func timestamptzToString(t pgtype.Timestamptz) *string {
	if !t.Valid {
		return nil
	}
	s := t.Time.Format(time.RFC3339)
	return &s
}

func int4ToIntPtr(i pgtype.Int4) *int {
	if !i.Valid {
		return nil
	}
	v := int(i.Int32)
	return &v
}

func buildMatchResponse(m repository.Match) MatchResponse {
	resp := MatchResponse{
		ID:             pgtypeUUIDToStringRequired(m.ID),
		Player1ID:      pgtypeUUIDToStringRequired(m.Player1ID),
		Player2ID:      pgtypeUUIDToStringRequired(m.Player2ID),
		Player1Partner: pgtypeUUIDToString(m.Player1PartnerID),
		Player2Partner: pgtypeUUIDToString(m.Player2PartnerID),
		EventID:        pgtypeUUIDToString(m.EventID),
		CommunityID:    pgtypeUUIDToString(m.CommunityID),
		Composition:    string(m.Composition),
		WinnerID:       pgtypeUUIDToString(m.WinnerID),
		SubmittedBy:    pgtypeUUIDToString(m.SubmittedBy),
		ConfirmedBy:    pgtypeUUIDToString(m.ConfirmedBy),
		SubmittedAt:    timestamptzToString(m.SubmittedAt),
		ConfirmedAt:    timestamptzToString(m.ConfirmedAt),
		RoundName:      pgtypeTextToStringPtr(m.RoundName),
		RoundNumber:    int4ToIntPtr(m.RoundNumber),
		CourtNumber:    int4ToIntPtr(m.CourtNumber),
		ScheduledTime:  timestamptzToString(m.ScheduledTime),
		PlayedAt:       timestamptzToString(m.PlayedAt),
		CreatedAt:      m.CreatedAt.Time.Format(time.RFC3339),
	}

	if m.ResultStatus.Valid {
		resp.ResultStatus = string(m.ResultStatus.ResultStatus)
	} else {
		resp.ResultStatus = "pending"
	}

	if m.DisputeReason.Valid {
		resp.DisputeReason = &m.DisputeReason.String
	}

	// Parse score from JSON
	if len(m.Score) > 0 {
		var score any
		if err := json.Unmarshal(m.Score, &score); err == nil {
			resp.Score = score
		}
	}

	// Include rating changes if available
	if m.Player1RatingBefore.Valid && m.Player1RatingAfter.Valid {
		p1Before := numericToFloat(m.Player1RatingBefore)
		p1After := numericToFloat(m.Player1RatingAfter)
		p2Before := numericToFloat(m.Player2RatingBefore)
		p2After := numericToFloat(m.Player2RatingAfter)
		resp.RatingChanges = &struct {
			Player1 *RatingChangeInfo `json:"player1,omitempty"`
			Player2 *RatingChangeInfo `json:"player2,omitempty"`
		}{
			Player1: &RatingChangeInfo{
				Before: p1Before,
				After:  p1After,
				Change: p1After - p1Before,
			},
			Player2: &RatingChangeInfo{
				Before: p2Before,
				After:  p2After,
				Change: p2After - p2Before,
			},
		}
	}

	return resp
}

// notifyMatchResultPending creates notifications for opponents when a result is submitted.
func (s *MatchService) notifyMatchResultPending(ctx context.Context, match repository.Match, submitterID uuid.UUID, matchID uuid.UUID) {
	if s.notifications == nil {
		return
	}

	playerIDs := []pgtype.UUID{
		match.Player1ID,
		match.Player2ID,
		match.Player1PartnerID,
		match.Player2PartnerID,
	}

	for _, pgID := range playerIDs {
		if !pgID.Valid {
			continue
		}
		id, err := uuid.FromBytes(pgID.Bytes[:])
		if err != nil || id == submitterID {
			continue
		}

		_, err = s.notifications.Create(ctx, id,
			"match_result_pending",
			"Подтвердите результат матча",
			"Ваш соперник внёс результат матча. Пожалуйста, подтвердите или оспорьте его.",
			map[string]any{
				"match_id": matchID.String(),
			},
		)
		if err != nil {
			slog.Warn("failed to create match_result_pending notification", "user_id", id, "error", err)
		}
	}
}

// notifyRatingChanged creates rating_changed notifications for winner and loser.
func (s *MatchService) notifyRatingChanged(ctx context.Context, winnerID, loserID uuid.UUID, matchID pgtype.UUID, change elo.RatingChange) {
	if s.notifications == nil {
		return
	}

	matchUUIDStr := ""
	if matchID.Valid {
		if id, err := uuid.FromBytes(matchID.Bytes[:]); err == nil {
			matchUUIDStr = id.String()
		}
	}

	if matchUUIDStr == "" {
		return
	}

	// Winner
	if _, err := s.notifications.Create(ctx, winnerID,
		"rating_changed",
		"Рейтинг обновлён",
		fmt.Sprintf("Ваш новый рейтинг: %.1f (%+.1f)", change.WinnerNewRating, change.WinnerDelta),
		map[string]any{
			"match_id": matchUUIDStr,
		},
	); err != nil {
		slog.Warn("failed to create winner rating_changed notification", "user_id", winnerID, "error", err)
	}

	// Loser
	if _, err := s.notifications.Create(ctx, loserID,
		"rating_changed",
		"Рейтинг обновлён",
		fmt.Sprintf("Ваш новый рейтинг: %.1f (%+.1f)", change.LoserNewRating, change.LoserDelta),
		map[string]any{
			"match_id": matchUUIDStr,
		},
	); err != nil {
		slog.Warn("failed to create loser rating_changed notification", "user_id", loserID, "error", err)
	}
}

func pgtypeTextToStringPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}
