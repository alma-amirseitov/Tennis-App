package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// CommunityService handles community business logic
type CommunityService struct {
	repo *repository.Queries
}

// NewCommunityService creates a new CommunityService
func NewCommunityService(repo *repository.Queries) *CommunityService {
	return &CommunityService{repo: repo}
}

// CreateCommunityInput represents input for creating a community
type CreateCommunityInput struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	CommunityType string `json:"community_type"`
	AccessLevel   string `json:"access_level"`
	District      string `json:"district"`
}

// Create creates a new community
func (s *CommunityService) Create(ctx context.Context, userID uuid.UUID, input CreateCommunityInput) (map[string]interface{}, error) {
	// Generate slug from name
	slug := generateSlug(input.Name)

	// Determine verification status
	verificationStatus := repository.NullVerificationStatus{
		VerificationStatus: repository.VerificationStatusNone, Valid: true,
	}
	if input.CommunityType == "club" || input.CommunityType == "league" || input.CommunityType == "organizer" {
		verificationStatus = repository.NullVerificationStatus{
			VerificationStatus: repository.VerificationStatusPending, Valid: true,
		}
	}

	accessLevel := repository.NullCommunityAccess{
		CommunityAccess: repository.CommunityAccessOpen, Valid: true,
	}
	if input.AccessLevel == "closed" {
		accessLevel = repository.NullCommunityAccess{
			CommunityAccess: repository.CommunityAccessClosed, Valid: true,
		}
	}

	community, err := s.repo.CreateCommunity(ctx, repository.CreateCommunityParams{
		Name:               input.Name,
		Slug:               pgtype.Text{String: slug, Valid: true},
		Description:        pgtype.Text{String: input.Description, Valid: true},
		CommunityType:      repository.CommunityType(input.CommunityType),
		AccessLevel:        accessLevel,
		VerificationStatus: verificationStatus,
		District:           pgtype.Text{String: input.District, Valid: true},
		CreatedBy:          pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("create community: %w", err)
	}

	// Add creator as owner
	_, err = s.repo.AddCommunityMember(ctx, repository.AddCommunityMemberParams{
		CommunityID: community.ID,
		UserID:      pgtype.UUID{Bytes: userID, Valid: true},
		Role:        repository.NullCommunityRole{CommunityRole: repository.CommunityRoleOwner, Valid: true},
		Status:      repository.NullMemberStatus{MemberStatus: repository.MemberStatusActive, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("add owner: %w", err)
	}

	return buildCommunityResponse(community), nil
}

// ListCommunitiesInput contains filter parameters for listing communities
type ListCommunitiesInput struct {
	Type         string
	AccessLevel  string
	VerifiedOnly bool
	District     string
	Query        string
	Sort         string
	Page         int
	PerPage      int
}

// List returns a paginated list of communities
func (s *CommunityService) List(ctx context.Context, input ListCommunitiesInput) ([]map[string]interface{}, *PaginationInfo, error) {
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PerPage < 1 || input.PerPage > 100 {
		input.PerPage = 20
	}
	if input.Sort == "" {
		input.Sort = "members"
	}

	offset := (input.Page - 1) * input.PerPage

	params := repository.ListCommunitiesParams{
		SortBy:       input.Sort,
		ResultLimit:  int32(input.PerPage),
		ResultOffset: int32(offset),
	}
	countParams := repository.CountCommunitiesParams{}

	if input.Type != "" {
		params.CommunityType = repository.NullCommunityType{CommunityType: repository.CommunityType(input.Type), Valid: true}
		countParams.CommunityType = params.CommunityType
	}
	if input.AccessLevel != "" {
		params.AccessLevel = repository.NullCommunityAccess{CommunityAccess: repository.CommunityAccess(input.AccessLevel), Valid: true}
		countParams.AccessLevel = params.AccessLevel
	}
	if input.VerifiedOnly {
		params.VerifiedOnly = pgtype.Bool{Bool: true, Valid: true}
		countParams.VerifiedOnly = params.VerifiedOnly
	}
	if input.District != "" {
		params.District = pgtype.Text{String: input.District, Valid: true}
		countParams.District = params.District
	}
	if input.Query != "" {
		params.Query = pgtype.Text{String: input.Query, Valid: true}
		countParams.Query = params.Query
	}

	communities, err := s.repo.ListCommunities(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("list communities: %w", err)
	}

	total, err := s.repo.CountCommunities(ctx, countParams)
	if err != nil {
		return nil, nil, fmt.Errorf("count communities: %w", err)
	}

	totalInt := int(total)
	totalPages := totalInt / input.PerPage
	if totalInt%input.PerPage > 0 {
		totalPages++
	}

	result := make([]map[string]interface{}, 0, len(communities))
	for _, c := range communities {
		cID, _ := uuid.FromBytes(c.ID.Bytes[:])
		result = append(result, map[string]interface{}{
			"id":                  cID.String(),
			"name":                c.Name,
			"slug":                c.Slug.String,
			"description":         c.Description.String,
			"community_type":      string(c.CommunityType),
			"access_level":        string(c.AccessLevel.CommunityAccess),
			"verification_status": string(c.VerificationStatus.VerificationStatus),
			"logo_url":            c.LogoUrl.String,
			"district":            c.District.String,
			"member_count":        c.MemberCount.Int32,
			"event_count":         c.EventCount.Int32,
		})
	}

	pagination := &PaginationInfo{
		Page:       input.Page,
		PerPage:    input.PerPage,
		Total:      totalInt,
		TotalPages: totalPages,
	}

	return result, pagination, nil
}

// GetByID returns community details with the current user's role
func (s *CommunityService) GetByID(ctx context.Context, userID, communityID uuid.UUID) (map[string]interface{}, error) {
	community, err := s.repo.GetCommunityByID(ctx, pgtype.UUID{Bytes: communityID, Valid: true})
	if err == pgx.ErrNoRows {
		return nil, ErrCommunityNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get community: %w", err)
	}

	result := buildCommunityResponse(community)

	// Get user's membership
	member, err := s.repo.GetCommunityMember(ctx, repository.GetCommunityMemberParams{
		CommunityID: pgtype.UUID{Bytes: communityID, Valid: true},
		UserID:      pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err == nil {
		result["my_role"] = string(member.Role.CommunityRole)
		result["my_status"] = string(member.Status.MemberStatus)
	} else {
		result["my_role"] = nil
		result["my_status"] = nil
	}

	return result, nil
}

// Join handles a user joining a community
func (s *CommunityService) Join(ctx context.Context, userID, communityID uuid.UUID, message string) (map[string]interface{}, error) {
	community, err := s.repo.GetCommunityByID(ctx, pgtype.UUID{Bytes: communityID, Valid: true})
	if err == pgx.ErrNoRows {
		return nil, ErrCommunityNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get community: %w", err)
	}

	// Check if already a member
	existing, err := s.repo.GetCommunityMember(ctx, repository.GetCommunityMemberParams{
		CommunityID: pgtype.UUID{Bytes: communityID, Valid: true},
		UserID:      pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err == nil {
		if existing.Status.MemberStatus == repository.MemberStatusActive {
			return nil, ErrAlreadyMember
		}
		if existing.Status.MemberStatus == repository.MemberStatusPending {
			return nil, ErrAlreadyMember.WithMessage("Join request already pending")
		}
		if existing.Status.MemberStatus == repository.MemberStatusBanned {
			return nil, ErrForbidden.WithMessage("You are banned from this community")
		}
	}

	// Determine status based on access level
	status := repository.NullMemberStatus{MemberStatus: repository.MemberStatusActive, Valid: true}
	responseMsg := "Вы вступили в сообщество"
	if community.AccessLevel.CommunityAccess == repository.CommunityAccessClosed {
		status = repository.NullMemberStatus{MemberStatus: repository.MemberStatusPending, Valid: true}
		responseMsg = "Заявка отправлена"
	}

	var appMsg pgtype.Text
	if message != "" {
		appMsg = pgtype.Text{String: message, Valid: true}
	}

	member, err := s.repo.AddCommunityMember(ctx, repository.AddCommunityMemberParams{
		CommunityID:        pgtype.UUID{Bytes: communityID, Valid: true},
		UserID:             pgtype.UUID{Bytes: userID, Valid: true},
		Role:               repository.NullCommunityRole{CommunityRole: repository.CommunityRoleMember, Valid: true},
		Status:             status,
		ApplicationMessage: appMsg,
	})
	if err != nil {
		return nil, fmt.Errorf("add member: %w", err)
	}

	return map[string]interface{}{
		"status":  string(member.Status.MemberStatus),
		"message": responseMsg,
	}, nil
}

// Leave handles a user leaving a community
func (s *CommunityService) Leave(ctx context.Context, userID, communityID uuid.UUID) error {
	member, err := s.repo.GetCommunityMember(ctx, repository.GetCommunityMemberParams{
		CommunityID: pgtype.UUID{Bytes: communityID, Valid: true},
		UserID:      pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err == pgx.ErrNoRows {
		return ErrNotCommunityMember
	}
	if err != nil {
		return fmt.Errorf("get member: %w", err)
	}

	// Owner cannot leave
	if member.Role.CommunityRole == repository.CommunityRoleOwner {
		return ErrForbidden.WithMessage("Owner cannot leave community. Transfer ownership first.")
	}

	err = s.repo.DeleteCommunityMember(ctx, repository.DeleteCommunityMemberParams{
		CommunityID: pgtype.UUID{Bytes: communityID, Valid: true},
		UserID:      pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("delete member: %w", err)
	}

	return nil
}

// ListMembersInput contains filter parameters for listing members
type ListMembersInput struct {
	CommunityID uuid.UUID
	Role        string
	Status      string
	Query       string
	Sort        string
	Page        int
	PerPage     int
}

// ListMembers returns members of a community
func (s *CommunityService) ListMembers(ctx context.Context, input ListMembersInput) ([]map[string]interface{}, *PaginationInfo, error) {
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PerPage < 1 || input.PerPage > 100 {
		input.PerPage = 20
	}
	if input.Sort == "" {
		input.Sort = "joined"
	}

	offset := (input.Page - 1) * input.PerPage

	params := repository.ListCommunityMembersParams{
		CommunityID:  pgtype.UUID{Bytes: input.CommunityID, Valid: true},
		SortBy:       input.Sort,
		ResultLimit:  int32(input.PerPage),
		ResultOffset: int32(offset),
	}
	countParams := repository.CountCommunityMembersParams{
		CommunityID: pgtype.UUID{Bytes: input.CommunityID, Valid: true},
	}

	if input.Role != "" {
		params.MemberRole = repository.NullCommunityRole{CommunityRole: repository.CommunityRole(input.Role), Valid: true}
		countParams.MemberRole = params.MemberRole
	}
	if input.Status != "" {
		params.MemberStatus = repository.NullMemberStatus{MemberStatus: repository.MemberStatus(input.Status), Valid: true}
		countParams.MemberStatus = params.MemberStatus
	}
	if input.Query != "" {
		params.Query = pgtype.Text{String: input.Query, Valid: true}
		countParams.Query = params.Query
	}

	members, err := s.repo.ListCommunityMembers(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("list members: %w", err)
	}

	total, err := s.repo.CountCommunityMembers(ctx, countParams)
	if err != nil {
		return nil, nil, fmt.Errorf("count members: %w", err)
	}

	totalInt := int(total)
	totalPages := totalInt / input.PerPage
	if totalInt%input.PerPage > 0 {
		totalPages++
	}

	result := make([]map[string]interface{}, 0, len(members))
	for _, m := range members {
		mUserID, _ := uuid.FromBytes(m.UserID.Bytes[:])
		result = append(result, map[string]interface{}{
			"user_id":          mUserID.String(),
			"first_name":       m.FirstName.String,
			"last_name":        m.LastName.String,
			"avatar_url":       m.AvatarUrl.String,
			"ntrp_level":       numericToFloat(m.NtrpLevel),
			"role":             string(m.Role.CommunityRole),
			"status":           string(m.Status.MemberStatus),
			"community_rating": numericToFloat(m.CommunityRating),
			"joined_at":        m.JoinedAt.Time,
		})
	}

	pagination := &PaginationInfo{
		Page:       input.Page,
		PerPage:    input.PerPage,
		Total:      totalInt,
		TotalPages: totalPages,
	}

	return result, pagination, nil
}

// UpdateMemberRole updates a member's role
func (s *CommunityService) UpdateMemberRole(ctx context.Context, actorID, communityID, targetUserID uuid.UUID, newRole string) error {
	// Verify actor has permission
	actor, err := s.repo.GetCommunityMember(ctx, repository.GetCommunityMemberParams{
		CommunityID: pgtype.UUID{Bytes: communityID, Valid: true},
		UserID:      pgtype.UUID{Bytes: actorID, Valid: true},
	})
	if err != nil {
		return ErrForbidden
	}

	// Only owner and admin can change roles
	if actor.Role.CommunityRole != repository.CommunityRoleOwner && actor.Role.CommunityRole != repository.CommunityRoleAdmin {
		return ErrInsufficientRole
	}

	// Only owner can promote to admin
	if newRole == "admin" && actor.Role.CommunityRole != repository.CommunityRoleOwner {
		return ErrInsufficientRole.WithMessage("Only owner can promote to admin")
	}

	// Cannot change owner's role
	target, err := s.repo.GetCommunityMember(ctx, repository.GetCommunityMemberParams{
		CommunityID: pgtype.UUID{Bytes: communityID, Valid: true},
		UserID:      pgtype.UUID{Bytes: targetUserID, Valid: true},
	})
	if err == pgx.ErrNoRows {
		return ErrNotCommunityMember
	}
	if err != nil {
		return fmt.Errorf("get target member: %w", err)
	}

	if target.Role.CommunityRole == repository.CommunityRoleOwner {
		return ErrForbidden.WithMessage("Cannot change owner's role")
	}

	_, err = s.repo.UpdateCommunityMemberRole(ctx, repository.UpdateCommunityMemberRoleParams{
		CommunityID: pgtype.UUID{Bytes: communityID, Valid: true},
		UserID:      pgtype.UUID{Bytes: targetUserID, Valid: true},
		Role:        repository.NullCommunityRole{CommunityRole: repository.CommunityRole(newRole), Valid: true},
	})
	if err != nil {
		return fmt.Errorf("update role: %w", err)
	}

	return nil
}

// ReviewRequest approves or rejects a join request
func (s *CommunityService) ReviewRequest(ctx context.Context, actorID, communityID, targetUserID uuid.UUID, approve bool) error {
	// Verify actor has permission
	actor, err := s.repo.GetCommunityMember(ctx, repository.GetCommunityMemberParams{
		CommunityID: pgtype.UUID{Bytes: communityID, Valid: true},
		UserID:      pgtype.UUID{Bytes: actorID, Valid: true},
	})
	if err != nil {
		return ErrForbidden
	}

	if actor.Role.CommunityRole != repository.CommunityRoleOwner && actor.Role.CommunityRole != repository.CommunityRoleAdmin && actor.Role.CommunityRole != repository.CommunityRoleModerator {
		return ErrInsufficientRole
	}

	// Check target is pending
	target, err := s.repo.GetCommunityMember(ctx, repository.GetCommunityMemberParams{
		CommunityID: pgtype.UUID{Bytes: communityID, Valid: true},
		UserID:      pgtype.UUID{Bytes: targetUserID, Valid: true},
	})
	if err == pgx.ErrNoRows {
		return ErrNotFound.WithMessage("Join request not found")
	}
	if err != nil {
		return fmt.Errorf("get target: %w", err)
	}

	if target.Status.MemberStatus != repository.MemberStatusPending {
		return ErrValidation.WithMessage("User is not in pending status")
	}

	newStatus := repository.NullMemberStatus{MemberStatus: repository.MemberStatusActive, Valid: true}
	if !approve {
		newStatus = repository.NullMemberStatus{MemberStatus: repository.MemberStatusLeft, Valid: true}
	}

	_, err = s.repo.UpdateCommunityMemberStatus(ctx, repository.UpdateCommunityMemberStatusParams{
		CommunityID: pgtype.UUID{Bytes: communityID, Valid: true},
		UserID:      pgtype.UUID{Bytes: targetUserID, Valid: true},
		Status:      newStatus,
		ReviewedBy:  pgtype.UUID{Bytes: actorID, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	return nil
}

// ListMyCommunities returns communities the user is a member of
func (s *CommunityService) ListMyCommunities(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error) {
	communities, err := s.repo.ListMyCommunities(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("list my communities: %w", err)
	}

	result := make([]map[string]interface{}, 0, len(communities))
	for _, c := range communities {
		cID, _ := uuid.FromBytes(c.ID.Bytes[:])
		result = append(result, map[string]interface{}{
			"id":                  cID.String(),
			"name":                c.Name,
			"slug":                c.Slug.String,
			"description":         c.Description.String,
			"community_type":      string(c.CommunityType),
			"access_level":        string(c.AccessLevel.CommunityAccess),
			"verification_status": string(c.VerificationStatus.VerificationStatus),
			"logo_url":            c.LogoUrl.String,
			"member_count":        c.MemberCount.Int32,
			"my_role":             string(c.Role.CommunityRole),
		})
	}

	return result, nil
}

func buildCommunityResponse(c repository.Community) map[string]interface{} {
	cID, _ := uuid.FromBytes(c.ID.Bytes[:])
	creatorID, _ := uuid.FromBytes(c.CreatedBy.Bytes[:])

	return map[string]interface{}{
		"id":                  cID.String(),
		"name":                c.Name,
		"slug":                c.Slug.String,
		"description":         c.Description.String,
		"rules":               c.Rules.String,
		"community_type":      string(c.CommunityType),
		"access_level":        string(c.AccessLevel.CommunityAccess),
		"verification_status": string(c.VerificationStatus.VerificationStatus),
		"logo_url":            c.LogoUrl.String,
		"banner_url":          c.BannerUrl.String,
		"contact_phone":       c.ContactPhone.String,
		"social_links":        c.SocialLinks,
		"address":             c.Address.String,
		"district":            c.District.String,
		"member_count":        c.MemberCount.Int32,
		"event_count":         c.EventCount.Int32,
		"created_by":          creatorID.String(),
		"created_at":          c.CreatedAt.Time,
	}
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove non-alphanumeric characters except hyphens
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	s := result.String()
	if s == "" {
		s = uuid.New().String()[:8]
	}
	return s
}
