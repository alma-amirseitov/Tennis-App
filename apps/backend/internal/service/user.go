package service

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// UserService handles user-related business logic
type UserService struct {
	repo    *repository.Queries
	storage *StorageService
}

// NewUserService creates a new UserService
func NewUserService(repo *repository.Queries, storage *StorageService) *UserService {
	return &UserService{
		repo:    repo,
		storage: storage,
	}
}

// GetProfile returns the full profile of the authenticated user
func (s *UserService) GetProfile(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	user, err := s.repo.GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	// Get stats
	stats, _ := s.repo.GetUserStats(ctx, pgtype.UUID{Bytes: userID, Valid: true})

	// Get badges
	badges, _ := s.repo.GetUserBadges(ctx, pgtype.UUID{Bytes: userID, Valid: true})

	// Get communities
	communities, _ := s.repo.GetUserCommunities(ctx, pgtype.UUID{Bytes: userID, Valid: true})

	profile := buildUserProfile(user, true)

	// Add stats
	if stats.UserID.Valid {
		profile["stats"] = map[string]interface{}{
			"total_games":    stats.TotalGames.Int32,
			"total_wins":     stats.TotalWins.Int32,
			"total_losses":   stats.TotalLosses.Int32,
			"win_rate":       numericToFloat(stats.WinRate),
			"current_streak": stats.CurrentStreak.Int32,
			"best_streak":    stats.BestStreak.Int32,
		}
	}

	// Add badges
	badgeList := make([]map[string]interface{}, 0, len(badges))
	for _, b := range badges {
		badgeList = append(badgeList, map[string]interface{}{
			"id":        b.BadgeID,
			"icon":      b.Icon.String,
			"name_ru":   b.NameRu,
			"name_kz":   b.NameKz.String,
			"name_en":   b.NameEn.String,
			"earned_at": b.EarnedAt.Time,
		})
	}
	profile["badges"] = badgeList

	// Add communities
	communityList := make([]map[string]interface{}, 0, len(communities))
	for _, c := range communities {
		cID, _ := uuid.FromBytes(c.ID.Bytes[:])
		communityList = append(communityList, map[string]interface{}{
			"id":   cID.String(),
			"name": c.Name,
			"slug": c.Slug.String,
			"role": string(c.Role.CommunityRole),
		})
	}
	profile["communities"] = communityList

	return profile, nil
}

// UpdateProfileInput represents fields that can be updated
type UpdateProfileInput struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Bio       *string `json:"bio"`
	District  *string `json:"district"`
	Language  *string `json:"language"`
	City      *string `json:"city"`
}

// UpdateProfile updates user profile fields
func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, input UpdateProfileInput) (map[string]interface{}, error) {
	params := repository.UpdateUserParams{
		ID: pgtype.UUID{Bytes: userID, Valid: true},
	}

	if input.FirstName != nil {
		params.FirstName = pgtype.Text{String: *input.FirstName, Valid: true}
	}
	if input.LastName != nil {
		params.LastName = pgtype.Text{String: *input.LastName, Valid: true}
	}
	if input.Bio != nil {
		params.Bio = pgtype.Text{String: *input.Bio, Valid: true}
	}
	if input.District != nil {
		params.District = pgtype.Text{String: *input.District, Valid: true}
	}
	if input.Language != nil {
		params.Language = pgtype.Text{String: *input.Language, Valid: true}
	}
	if input.City != nil {
		params.City = pgtype.Text{String: *input.City, Valid: true}
	}

	user, err := s.repo.UpdateUser(ctx, params)
	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return buildUserProfile(user, true), nil
}

// GetPublicProfile returns a public view of a user's profile
func (s *UserService) GetPublicProfile(ctx context.Context, currentUserID, targetUserID uuid.UUID) (map[string]interface{}, error) {
	user, err := s.repo.GetUserByID(ctx, pgtype.UUID{Bytes: targetUserID, Valid: true})
	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	profile := buildUserProfile(user, false)

	// Get stats if visible
	if user.ShowStats.Bool {
		stats, _ := s.repo.GetUserStats(ctx, pgtype.UUID{Bytes: targetUserID, Valid: true})
		if stats.UserID.Valid {
			profile["stats"] = map[string]interface{}{
				"total_games":    stats.TotalGames.Int32,
				"total_wins":     stats.TotalWins.Int32,
				"total_losses":   stats.TotalLosses.Int32,
				"win_rate":       numericToFloat(stats.WinRate),
				"current_streak": stats.CurrentStreak.Int32,
			}
		}
	}

	// Get badges
	badges, _ := s.repo.GetUserBadges(ctx, pgtype.UUID{Bytes: targetUserID, Valid: true})
	badgeList := make([]map[string]interface{}, 0, len(badges))
	for _, b := range badges {
		badgeList = append(badgeList, map[string]interface{}{
			"id":        b.BadgeID,
			"icon":      b.Icon.String,
			"name_ru":   b.NameRu,
			"earned_at": b.EarnedAt.Time,
		})
	}
	profile["badges"] = badgeList

	// Get communities
	communities, _ := s.repo.GetUserCommunities(ctx, pgtype.UUID{Bytes: targetUserID, Valid: true})
	communityList := make([]map[string]interface{}, 0, len(communities))
	for _, c := range communities {
		cID, _ := uuid.FromBytes(c.ID.Bytes[:])
		communityList = append(communityList, map[string]interface{}{
			"id":   cID.String(),
			"name": c.Name,
			"role": string(c.Role.CommunityRole),
		})
	}
	profile["communities"] = communityList

	// Check friendship
	isFriend, _ := s.repo.CheckFriendship(ctx, repository.CheckFriendshipParams{
		UserID:   pgtype.UUID{Bytes: currentUserID, Valid: true},
		FriendID: pgtype.UUID{Bytes: targetUserID, Valid: true},
	})
	profile["is_friend"] = isFriend

	// Count mutual communities
	mutualCount, _ := s.repo.CountMutualCommunities(ctx, repository.CountMutualCommunitiesParams{
		UserID:   pgtype.UUID{Bytes: currentUserID, Valid: true},
		UserID_2: pgtype.UUID{Bytes: targetUserID, Valid: true},
	})
	profile["mutual_communities"] = mutualCount

	return profile, nil
}

// SearchUsersInput contains search parameters
type SearchUsersInput struct {
	Query    string
	MinLevel *float64
	MaxLevel *float64
	Gender   string
	District string
	Sort     string
	Page     int
	PerPage  int
}

// SearchUsersResult contains search results with pagination
type SearchUsersResult struct {
	Users      []map[string]interface{} `json:"users"`
	Pagination PaginationInfo           `json:"pagination"`
}

// PaginationInfo contains pagination metadata
type PaginationInfo struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// SearchUsers searches for users with filters
func (s *UserService) SearchUsers(ctx context.Context, input SearchUsersInput) (*SearchUsersResult, error) {
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PerPage < 1 || input.PerPage > 50 {
		input.PerPage = 20
	}
	if input.Sort == "" {
		input.Sort = "rating"
	}

	offset := (input.Page - 1) * input.PerPage

	searchParams := repository.SearchUsersParams{
		SortBy:       input.Sort,
		ResultLimit:  int32(input.PerPage),
		ResultOffset: int32(offset),
	}
	countParams := repository.CountSearchUsersParams{}

	if input.Query != "" {
		searchParams.Query = pgtype.Text{String: input.Query, Valid: true}
		countParams.Query = pgtype.Text{String: input.Query, Valid: true}
	}
	if input.MinLevel != nil {
		searchParams.MinLevel = pgtype.Numeric{Valid: true}
		searchParams.MinLevel.Scan(fmt.Sprintf("%.1f", *input.MinLevel))
		countParams.MinLevel = searchParams.MinLevel
	}
	if input.MaxLevel != nil {
		searchParams.MaxLevel = pgtype.Numeric{Valid: true}
		searchParams.MaxLevel.Scan(fmt.Sprintf("%.1f", *input.MaxLevel))
		countParams.MaxLevel = searchParams.MaxLevel
	}
	if input.District != "" {
		searchParams.District = pgtype.Text{String: input.District, Valid: true}
		countParams.District = pgtype.Text{String: input.District, Valid: true}
	}
	if input.Gender != "" {
		searchParams.Gender = repository.NullGenderType{GenderType: repository.GenderType(input.Gender), Valid: true}
		countParams.Gender = searchParams.Gender
	}

	users, err := s.repo.SearchUsers(ctx, searchParams)
	if err != nil {
		return nil, fmt.Errorf("search users: %w", err)
	}

	total, err := s.repo.CountSearchUsers(ctx, countParams)
	if err != nil {
		return nil, fmt.Errorf("count users: %w", err)
	}

	totalInt := int(total)
	totalPages := totalInt / input.PerPage
	if totalInt%input.PerPage > 0 {
		totalPages++
	}

	userList := make([]map[string]interface{}, 0, len(users))
	for _, u := range users {
		uid, _ := uuid.FromBytes(u.ID.Bytes[:])
		userList = append(userList, map[string]interface{}{
			"id":            uid.String(),
			"first_name":    u.FirstName.String,
			"last_name":     u.LastName.String,
			"avatar_url":    u.AvatarUrl.String,
			"ntrp_level":    numericToFloat(u.NtrpLevel),
			"level_label":   u.LevelLabel.String,
			"global_rating": numericToFloat(u.GlobalRating),
			"district":      u.District.String,
		})
	}

	return &SearchUsersResult{
		Users: userList,
		Pagination: PaginationInfo{
			Page:       input.Page,
			PerPage:    input.PerPage,
			Total:      totalInt,
			TotalPages: totalPages,
		},
	}, nil
}

// UploadAvatar handles avatar upload
func (s *UserService) UploadAvatar(ctx context.Context, userID uuid.UUID, fileData []byte) (string, error) {
	// Validate image
	contentType, err := ValidateImage(fileData)
	if err != nil {
		return "", err
	}

	// Generate key
	key := GenerateAvatarKey(userID, contentType)

	// Upload to S3
	url, err := s.storage.Upload(ctx, "", key, bytes.NewReader(fileData), contentType)
	if err != nil {
		return "", fmt.Errorf("upload avatar: %w", err)
	}

	// Update user avatar_url
	_, err = s.repo.UpdateUserAvatarURL(ctx, repository.UpdateUserAvatarURLParams{
		ID:        pgtype.UUID{Bytes: userID, Valid: true},
		AvatarUrl: pgtype.Text{String: url, Valid: true},
	})
	if err != nil {
		slog.Error("failed to update avatar URL in DB", "error", err, "user_id", userID)
		return "", fmt.Errorf("update avatar url: %w", err)
	}

	return url, nil
}

// buildUserProfile constructs a user profile map
func buildUserProfile(user repository.User, includePrivate bool) map[string]interface{} {
	uid, _ := uuid.FromBytes(user.ID.Bytes[:])

	profile := map[string]interface{}{
		"id":                  uid.String(),
		"first_name":          user.FirstName.String,
		"last_name":           user.LastName.String,
		"avatar_url":          user.AvatarUrl.String,
		"ntrp_level":          numericToFloat(user.NtrpLevel),
		"level_label":         user.LevelLabel.String,
		"global_rating":       numericToFloat(user.GlobalRating),
		"global_games_count":  user.GlobalGamesCount.Int32,
		"district":            user.District.String,
		"city":                user.City.String,
		"bio":                 user.Bio.String,
		"is_profile_complete": user.IsProfileComplete.Bool,
		"created_at":          user.CreatedAt.Time,
	}

	if user.Gender.Valid {
		profile["gender"] = string(user.Gender.GenderType)
	}
	if user.BirthYear.Valid {
		profile["birth_year"] = user.BirthYear.Int16
	}

	if includePrivate {
		profile["phone"] = maskPhone(user.Phone)
		profile["language"] = user.Language.String
		profile["quiz_completed"] = user.QuizCompleted.Bool
		profile["pin_set"] = user.PinHash.Valid && user.PinHash.String != ""
		profile["notification_settings"] = user.NotificationSettings
		profile["profile_visibility"] = user.ProfileVisibility.String
		profile["allow_messages_from"] = user.AllowMessagesFrom.String
		profile["show_stats"] = user.ShowStats.Bool
	} else {
		// For public profile: mask last name
		if user.LastName.Valid && len(user.LastName.String) > 0 {
			profile["last_name"] = string([]rune(user.LastName.String)[0:1]) + "."
		}
	}

	return profile
}

// numericToFloat converts pgtype.Numeric to float64
func numericToFloat(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, _ := n.Float64Value()
	return f.Float64
}

