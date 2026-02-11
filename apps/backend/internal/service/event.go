package service

import (
	"context"
	"fmt"
	"time"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// EventService handles event business logic
type EventService struct {
	repo *repository.Queries
}

// NewEventService creates a new EventService
func NewEventService(repo *repository.Queries) *EventService {
	return &EventService{repo: repo}
}

// CreateEventInput represents input for creating an event
type CreateEventInput struct {
	Title              string     `json:"title"`
	Description        string     `json:"description"`
	EventType          string     `json:"event_type"`
	CommunityID        *string    `json:"community_id"`
	PlayerComposition  string     `json:"player_composition"`
	MatchFormat        string     `json:"match_format"`
	MatchFormatDetails any        `json:"match_format_details"`
	CourtID            *string    `json:"court_id"`
	LocationName       string     `json:"location_name"`
	LocationAddress    string     `json:"location_address"`
	StartTime          time.Time  `json:"start_time"`
	EndTime            *time.Time `json:"end_time"`
	MaxParticipants    int        `json:"max_participants"`
	MinParticipants    int        `json:"min_participants"`
	MinLevel           *float64   `json:"min_level"`
	MaxLevel           *float64   `json:"max_level"`
	GenderRestriction  string     `json:"gender_restriction"`
	RegistrationDeadline *time.Time `json:"registration_deadline"`
	IsPaid             bool       `json:"is_paid"`
	PriceAmount        *float64   `json:"price_amount"`
	Status             string     `json:"status"`
}

// Create creates a new event
func (s *EventService) Create(ctx context.Context, userID uuid.UUID, input CreateEventInput) (map[string]interface{}, error) {
	if input.MinParticipants < 2 {
		input.MinParticipants = 2
	}
	if input.MaxParticipants < input.MinParticipants {
		return nil, ErrValidation.WithMessage("max_participants must be >= min_participants")
	}
	if input.MinLevel != nil && input.MaxLevel != nil && *input.MinLevel > *input.MaxLevel {
		return nil, ErrValidation.WithMessage("min_level must be <= max_level")
	}
	if input.StartTime.Before(time.Now()) {
		return nil, ErrValidation.WithMessage("start_time must be in the future")
	}

	status := repository.NullEventStatus{EventStatus: repository.EventStatusPublished, Valid: true}
	if input.Status == "draft" {
		status = repository.NullEventStatus{EventStatus: repository.EventStatusDraft, Valid: true}
	}

	composition := repository.PlayerCompositionSingles
	if input.PlayerComposition != "" {
		composition = repository.PlayerComposition(input.PlayerComposition)
	}

	matchFormat := repository.NullMatchFormat{MatchFormat: repository.MatchFormatBestOf, Valid: true}
	if input.MatchFormat != "" {
		matchFormat = repository.NullMatchFormat{MatchFormat: repository.MatchFormat(input.MatchFormat), Valid: true}
	}

	params := repository.CreateEventParams{
		Title:              input.Title,
		Description:        pgtype.Text{String: input.Description, Valid: input.Description != ""},
		EventType:          repository.EventType(input.EventType),
		Status:             status,
		PlayerComposition:  composition,
		MatchFormat:        matchFormat,
		LocationName:       pgtype.Text{String: input.LocationName, Valid: input.LocationName != ""},
		LocationAddress:    pgtype.Text{String: input.LocationAddress, Valid: input.LocationAddress != ""},
		StartTime:          pgtype.Timestamptz{Time: input.StartTime, Valid: true},
		MaxParticipants:    pgtype.Int4{Int32: int32(input.MaxParticipants), Valid: true},
		MinParticipants:    pgtype.Int4{Int32: int32(input.MinParticipants), Valid: true},
		IsPaid:             pgtype.Bool{Bool: input.IsPaid, Valid: true},
		CreatedBy:          pgtype.UUID{Bytes: userID, Valid: true},
	}

	if input.CommunityID != nil {
		cID, err := uuid.Parse(*input.CommunityID)
		if err == nil {
			params.CommunityID = pgtype.UUID{Bytes: cID, Valid: true}
		}
	}
	if input.CourtID != nil {
		courtID, err := uuid.Parse(*input.CourtID)
		if err == nil {
			params.CourtID = pgtype.UUID{Bytes: courtID, Valid: true}
		}
	}
	if input.EndTime != nil {
		params.EndTime = pgtype.Timestamptz{Time: *input.EndTime, Valid: true}
	}
	if input.MinLevel != nil {
		params.MinLevel = pgtype.Numeric{Valid: true}
		params.MinLevel.Scan(fmt.Sprintf("%.1f", *input.MinLevel))
	}
	if input.MaxLevel != nil {
		params.MaxLevel = pgtype.Numeric{Valid: true}
		params.MaxLevel.Scan(fmt.Sprintf("%.1f", *input.MaxLevel))
	}
	if input.GenderRestriction != "" {
		params.GenderRestriction = repository.NullGenderType{GenderType: repository.GenderType(input.GenderRestriction), Valid: true}
	}
	if input.RegistrationDeadline != nil {
		params.RegistrationDeadline = pgtype.Timestamptz{Time: *input.RegistrationDeadline, Valid: true}
	}
	if input.PriceAmount != nil {
		params.PriceAmount = pgtype.Numeric{Valid: true}
		params.PriceAmount.Scan(fmt.Sprintf("%.2f", *input.PriceAmount))
	}

	event, err := s.repo.CreateEvent(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("create event: %w", err)
	}

	// Auto-join creator as participant
	_, _ = s.repo.AddEventParticipant(ctx, repository.AddEventParticipantParams{
		EventID: event.ID,
		UserID:  pgtype.UUID{Bytes: userID, Valid: true},
		Status:  repository.NullParticipantStatus{ParticipantStatus: repository.ParticipantStatusRegistered, Valid: true},
	})

	return buildEventResponse(event), nil
}

// ListEventsInput contains filter parameters
type ListEventsInput struct {
	EventType   string
	Status      string
	Composition string
	CommunityID string
	MinLevel    *float64
	MaxLevel    *float64
	DateFrom    *time.Time
	DateTo      *time.Time
	District    string
	Sort        string
	Page        int
	PerPage     int
}

// List returns a paginated list of events
func (s *EventService) List(ctx context.Context, input ListEventsInput) ([]map[string]interface{}, *PaginationInfo, error) {
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PerPage < 1 || input.PerPage > 100 {
		input.PerPage = 20
	}
	if input.Sort == "" {
		input.Sort = "date_desc"
	}

	offset := (input.Page - 1) * input.PerPage

	params := repository.ListEventsParams{
		SortBy:       input.Sort,
		ResultLimit:  int32(input.PerPage),
		ResultOffset: int32(offset),
	}
	countParams := repository.CountEventsParams{}

	if input.EventType != "" {
		params.EventType = repository.NullEventType{EventType: repository.EventType(input.EventType), Valid: true}
		countParams.EventType = params.EventType
	}
	if input.Status != "" {
		params.EventStatus = repository.NullEventStatus{EventStatus: repository.EventStatus(input.Status), Valid: true}
		countParams.EventStatus = params.EventStatus
	}
	if input.Composition != "" {
		params.Composition = repository.NullPlayerComposition{PlayerComposition: repository.PlayerComposition(input.Composition), Valid: true}
		countParams.Composition = params.Composition
	}
	if input.CommunityID != "" {
		cID, err := uuid.Parse(input.CommunityID)
		if err == nil {
			params.CommunityID = pgtype.UUID{Bytes: cID, Valid: true}
			countParams.CommunityID = params.CommunityID
		}
	}
	if input.MinLevel != nil {
		params.MinLevel = pgtype.Numeric{Valid: true}
		params.MinLevel.Scan(fmt.Sprintf("%.1f", *input.MinLevel))
		countParams.MinLevel = params.MinLevel
	}
	if input.MaxLevel != nil {
		params.MaxLevel = pgtype.Numeric{Valid: true}
		params.MaxLevel.Scan(fmt.Sprintf("%.1f", *input.MaxLevel))
		countParams.MaxLevel = params.MaxLevel
	}
	if input.DateFrom != nil {
		params.DateFrom = pgtype.Timestamptz{Time: *input.DateFrom, Valid: true}
		countParams.DateFrom = params.DateFrom
	}
	if input.DateTo != nil {
		params.DateTo = pgtype.Timestamptz{Time: *input.DateTo, Valid: true}
		countParams.DateTo = params.DateTo
	}
	if input.District != "" {
		params.District = pgtype.Text{String: input.District, Valid: true}
		countParams.District = params.District
	}

	events, err := s.repo.ListEvents(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("list events: %w", err)
	}

	total, err := s.repo.CountEvents(ctx, countParams)
	if err != nil {
		return nil, nil, fmt.Errorf("count events: %w", err)
	}

	totalInt := int(total)
	totalPages := totalInt / input.PerPage
	if totalInt%input.PerPage > 0 {
		totalPages++
	}

	result := make([]map[string]interface{}, 0, len(events))
	for _, e := range events {
		result = append(result, buildEventListItem(e))
	}

	pagination := &PaginationInfo{
		Page:       input.Page,
		PerPage:    input.PerPage,
		Total:      totalInt,
		TotalPages: totalPages,
	}

	return result, pagination, nil
}

// GetByID returns event details with participant info
func (s *EventService) GetByID(ctx context.Context, userID, eventID uuid.UUID) (map[string]interface{}, error) {
	event, err := s.repo.GetEventByID(ctx, pgtype.UUID{Bytes: eventID, Valid: true})
	if err == pgx.ErrNoRows {
		return nil, ErrEventNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get event: %w", err)
	}

	result := buildEventResponse(event)

	// Get participants
	participants, _ := s.repo.ListEventParticipants(ctx, pgtype.UUID{Bytes: eventID, Valid: true})
	participantList := make([]map[string]interface{}, 0, len(participants))
	for _, p := range participants {
		pUID, _ := uuid.FromBytes(p.UserID.Bytes[:])
		participantList = append(participantList, map[string]interface{}{
			"id":         pUID.String(),
			"first_name": p.FirstName.String,
			"last_name":  p.LastName.String,
			"avatar_url": p.AvatarUrl.String,
			"ntrp_level": numericToFloat(p.NtrpLevel),
			"status":     string(p.Status.ParticipantStatus),
		})
	}
	result["participants"] = participantList

	// Check user's status
	myParticipant, err := s.repo.GetEventParticipant(ctx, repository.GetEventParticipantParams{
		EventID: pgtype.UUID{Bytes: eventID, Valid: true},
		UserID:  pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err == nil {
		result["my_status"] = string(myParticipant.Status.ParticipantStatus)
	} else {
		result["my_status"] = nil
	}

	// Check if user can join
	creatorID, _ := uuid.FromBytes(event.CreatedBy.Bytes[:])
	canJoin := result["my_status"] == nil &&
		(event.Status.EventStatus == repository.EventStatusPublished || event.Status.EventStatus == repository.EventStatusRegistrationOpen) &&
		event.CurrentParticipants.Int32 < event.MaxParticipants.Int32
	result["can_join"] = canJoin
	result["can_edit"] = userID == creatorID

	return result, nil
}

// Join adds a user to an event
func (s *EventService) Join(ctx context.Context, userID, eventID uuid.UUID) (map[string]interface{}, error) {
	event, err := s.repo.GetEventByID(ctx, pgtype.UUID{Bytes: eventID, Valid: true})
	if err == pgx.ErrNoRows {
		return nil, ErrEventNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get event: %w", err)
	}

	// Check status
	if event.Status.EventStatus != repository.EventStatusPublished && event.Status.EventStatus != repository.EventStatusRegistrationOpen {
		return nil, ErrEventClosed
	}

	// Check capacity
	if event.MaxParticipants.Valid && event.CurrentParticipants.Int32 >= event.MaxParticipants.Int32 {
		return nil, ErrEventFull
	}

	// Check registration deadline
	if event.RegistrationDeadline.Valid && time.Now().After(event.RegistrationDeadline.Time) {
		return nil, ErrEventClosed.WithMessage("Registration deadline passed")
	}

	// Check level
	user, err := s.repo.GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	userLevel := numericToFloat(user.NtrpLevel)
	if event.MinLevel.Valid && userLevel > 0 {
		minLevel := numericToFloat(event.MinLevel)
		maxLevel := numericToFloat(event.MaxLevel)
		if minLevel > 0 && userLevel < minLevel {
			return nil, ErrEventWrongLevel
		}
		if maxLevel > 0 && userLevel > maxLevel {
			return nil, ErrEventWrongLevel
		}
	}

	// Check not already joined
	_, err = s.repo.GetEventParticipant(ctx, repository.GetEventParticipantParams{
		EventID: pgtype.UUID{Bytes: eventID, Valid: true},
		UserID:  pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err == nil {
		return nil, ErrAlreadyJoinedEvent
	}

	participant, err := s.repo.AddEventParticipant(ctx, repository.AddEventParticipantParams{
		EventID: pgtype.UUID{Bytes: eventID, Valid: true},
		UserID:  pgtype.UUID{Bytes: userID, Valid: true},
		Status:  repository.NullParticipantStatus{ParticipantStatus: repository.ParticipantStatusRegistered, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("add participant: %w", err)
	}

	pID, _ := uuid.FromBytes(participant.ID.Bytes[:])

	return map[string]interface{}{
		"participant_id": pID.String(),
		"status":         string(participant.Status.ParticipantStatus),
	}, nil
}

// Leave removes a user from an event
func (s *EventService) Leave(ctx context.Context, userID, eventID uuid.UUID) error {
	event, err := s.repo.GetEventByID(ctx, pgtype.UUID{Bytes: eventID, Valid: true})
	if err == pgx.ErrNoRows {
		return ErrEventNotFound
	}
	if err != nil {
		return fmt.Errorf("get event: %w", err)
	}

	// Can only leave open/filling events
	if event.Status.EventStatus != repository.EventStatusPublished && event.Status.EventStatus != repository.EventStatusRegistrationOpen {
		return ErrValidation.WithMessage("Cannot leave event with status: " + string(event.Status.EventStatus))
	}

	// Check if participant exists
	_, err = s.repo.GetEventParticipant(ctx, repository.GetEventParticipantParams{
		EventID: pgtype.UUID{Bytes: eventID, Valid: true},
		UserID:  pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err == pgx.ErrNoRows {
		return ErrNotFound.WithMessage("You are not a participant")
	}
	if err != nil {
		return fmt.Errorf("get participant: %w", err)
	}

	err = s.repo.RemoveEventParticipant(ctx, repository.RemoveEventParticipantParams{
		EventID: pgtype.UUID{Bytes: eventID, Valid: true},
		UserID:  pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("remove participant: %w", err)
	}

	return nil
}

// Valid status transitions
var validTransitions = map[repository.EventStatus][]repository.EventStatus{
	repository.EventStatusDraft:              {repository.EventStatusPublished},
	repository.EventStatusPublished:          {repository.EventStatusRegistrationOpen, repository.EventStatusCancelled},
	repository.EventStatusRegistrationOpen:   {repository.EventStatusRegistrationClosed, repository.EventStatusCancelled},
	repository.EventStatusRegistrationClosed: {repository.EventStatusInProgress, repository.EventStatusCancelled},
	repository.EventStatusInProgress:         {repository.EventStatusCompleted, repository.EventStatusCancelled},
}

// UpdateStatus changes the event status with lifecycle validation
func (s *EventService) UpdateStatus(ctx context.Context, userID, eventID uuid.UUID, newStatus string) (map[string]interface{}, error) {
	event, err := s.repo.GetEventByID(ctx, pgtype.UUID{Bytes: eventID, Valid: true})
	if err == pgx.ErrNoRows {
		return nil, ErrEventNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get event: %w", err)
	}

	// Check permission — only creator can change status
	creatorID, _ := uuid.FromBytes(event.CreatedBy.Bytes[:])
	if userID != creatorID {
		return nil, ErrForbidden.WithMessage("Only event creator can change status")
	}

	targetStatus := repository.EventStatus(newStatus)
	currentStatus := event.Status.EventStatus

	// Cancelled allowed from any non-completed status
	if targetStatus == repository.EventStatusCancelled {
		if currentStatus == repository.EventStatusCompleted {
			return nil, ErrValidation.WithMessage("Cannot cancel completed event")
		}
	} else {
		// Validate transition
		allowed, ok := validTransitions[currentStatus]
		if !ok {
			return nil, ErrValidation.WithMessage("Invalid current status: " + string(currentStatus))
		}
		valid := false
		for _, s := range allowed {
			if s == targetStatus {
				valid = true
				break
			}
		}
		if !valid {
			return nil, ErrValidation.WithMessage(
				fmt.Sprintf("Cannot transition from %s to %s", currentStatus, targetStatus),
			)
		}
	}

	updated, err := s.repo.UpdateEventStatus(ctx, repository.UpdateEventStatusParams{
		ID:     pgtype.UUID{Bytes: eventID, Valid: true},
		Status: repository.NullEventStatus{EventStatus: targetStatus, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("update status: %w", err)
	}

	return buildUpdateEventStatusResponse(updated), nil
}

// GetCalendar returns events grouped by day for a given month
func (s *EventService) GetCalendar(ctx context.Context, year, month int, communityID string) (map[string][]map[string]interface{}, error) {
	monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, 0)

	params := repository.GetCalendarEventsParams{
		MonthStart: pgtype.Timestamptz{Time: monthStart, Valid: true},
		MonthEnd:   pgtype.Timestamptz{Time: monthEnd, Valid: true},
	}
	if communityID != "" {
		cID, err := uuid.Parse(communityID)
		if err == nil {
			params.CommunityID = pgtype.UUID{Bytes: cID, Valid: true}
		}
	}

	events, err := s.repo.GetCalendarEvents(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("get calendar events: %w", err)
	}

	result := make(map[string][]map[string]interface{})
	for _, e := range events {
		eID, _ := uuid.FromBytes(e.ID.Bytes[:])
		day := e.StartTime.Time.Format("2006-01-02")
		result[day] = append(result[day], map[string]interface{}{
			"id":                   eID.String(),
			"title":                e.Title,
			"event_type":           string(e.EventType),
			"status":               string(e.Status.EventStatus),
			"start_time":           e.StartTime.Time,
			"current_participants": e.CurrentParticipants.Int32,
			"max_participants":     e.MaxParticipants.Int32,
		})
	}

	return result, nil
}

// GetMyEvents returns events for the current user
func (s *EventService) GetMyEvents(ctx context.Context, userID uuid.UUID, tab string, page, perPage int) ([]map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	userPgID := pgtype.UUID{Bytes: userID, Valid: true}

	var result []map[string]interface{}

	switch tab {
	case "created":
		events, err := s.repo.ListMyCreatedEvents(ctx, repository.ListMyCreatedEventsParams{
			CreatedBy: userPgID,
			Limit:     int32(perPage),
			Offset:    int32(offset),
		})
		if err != nil {
			return nil, fmt.Errorf("list created events: %w", err)
		}
		result = make([]map[string]interface{}, 0, len(events))
		for _, e := range events {
			eID, _ := uuid.FromBytes(e.ID.Bytes[:])
			result = append(result, map[string]interface{}{
				"id":                   eID.String(),
				"title":                e.Title,
				"event_type":           string(e.EventType),
				"status":               string(e.Status.EventStatus),
				"start_time":           e.StartTime.Time,
				"current_participants": e.CurrentParticipants.Int32,
				"max_participants":     e.MaxParticipants.Int32,
			})
		}
	case "past":
		events, err := s.repo.ListMyPastEvents(ctx, repository.ListMyPastEventsParams{
			UserID: userPgID,
			Limit:  int32(perPage),
			Offset: int32(offset),
		})
		if err != nil {
			return nil, fmt.Errorf("list past events: %w", err)
		}
		result = make([]map[string]interface{}, 0, len(events))
		for _, e := range events {
			eID, _ := uuid.FromBytes(e.ID.Bytes[:])
			result = append(result, map[string]interface{}{
				"id":                   eID.String(),
				"title":                e.Title,
				"event_type":           string(e.EventType),
				"status":               string(e.Status.EventStatus),
				"start_time":           e.StartTime.Time,
				"current_participants": e.CurrentParticipants.Int32,
				"max_participants":     e.MaxParticipants.Int32,
			})
		}
	default: // "joined"
		events, err := s.repo.ListMyJoinedEvents(ctx, repository.ListMyJoinedEventsParams{
			UserID: userPgID,
			Limit:  int32(perPage),
			Offset: int32(offset),
		})
		if err != nil {
			return nil, fmt.Errorf("list joined events: %w", err)
		}
		result = make([]map[string]interface{}, 0, len(events))
		for _, e := range events {
			eID, _ := uuid.FromBytes(e.ID.Bytes[:])
			result = append(result, map[string]interface{}{
				"id":                   eID.String(),
				"title":                e.Title,
				"event_type":           string(e.EventType),
				"status":               string(e.Status.EventStatus),
				"start_time":           e.StartTime.Time,
				"current_participants": e.CurrentParticipants.Int32,
				"max_participants":     e.MaxParticipants.Int32,
			})
		}
	}

	return result, nil
}

func buildEventResponse(e repository.Event) map[string]interface{} {
	eID, _ := uuid.FromBytes(e.ID.Bytes[:])
	creatorID, _ := uuid.FromBytes(e.CreatedBy.Bytes[:])

	result := map[string]interface{}{
		"id":                   eID.String(),
		"title":                e.Title,
		"description":          e.Description.String,
		"event_type":           string(e.EventType),
		"status":               string(e.Status.EventStatus),
		"player_composition":   string(e.PlayerComposition),
		"location_name":        e.LocationName.String,
		"location_address":     e.LocationAddress.String,
		"start_time":           e.StartTime.Time,
		"max_participants":     e.MaxParticipants.Int32,
		"min_participants":     e.MinParticipants.Int32,
		"current_participants": e.CurrentParticipants.Int32,
		"min_level":            numericToFloat(e.MinLevel),
		"max_level":            numericToFloat(e.MaxLevel),
		"is_paid":              e.IsPaid.Bool,
		"created_by":           creatorID.String(),
		"created_at":           e.CreatedAt.Time,
	}

	if e.CommunityID.Valid {
		cID, _ := uuid.FromBytes(e.CommunityID.Bytes[:])
		result["community_id"] = cID.String()
	}
	if e.CourtID.Valid {
		courtID, _ := uuid.FromBytes(e.CourtID.Bytes[:])
		result["court_id"] = courtID.String()
	}
	if e.EndTime.Valid {
		result["end_time"] = e.EndTime.Time
	}
	if e.MatchFormat.Valid {
		result["match_format"] = string(e.MatchFormat.MatchFormat)
	}
	if e.RegistrationDeadline.Valid {
		result["registration_deadline"] = e.RegistrationDeadline.Time
	}
	if e.GenderRestriction.Valid {
		result["gender_restriction"] = string(e.GenderRestriction.GenderType)
	}

	return result
}

func buildUpdateEventStatusResponse(e repository.UpdateEventStatusRow) map[string]interface{} {
	eID, _ := uuid.FromBytes(e.ID.Bytes[:])
	creatorID, _ := uuid.FromBytes(e.CreatedBy.Bytes[:])

	result := map[string]interface{}{
		"id":                   eID.String(),
		"title":                e.Title,
		"description":          e.Description.String,
		"event_type":           string(e.EventType),
		"status":               string(e.Status.EventStatus),
		"player_composition":   string(e.PlayerComposition),
		"location_name":        e.LocationName.String,
		"location_address":     e.LocationAddress.String,
		"start_time":           e.StartTime.Time,
		"max_participants":     e.MaxParticipants.Int32,
		"min_participants":     e.MinParticipants.Int32,
		"current_participants": e.CurrentParticipants.Int32,
		"min_level":            numericToFloat(e.MinLevel),
		"max_level":            numericToFloat(e.MaxLevel),
		"is_paid":              e.IsPaid.Bool,
		"created_by":           creatorID.String(),
		"created_at":           e.CreatedAt.Time,
	}

	if e.CommunityID.Valid {
		cID, _ := uuid.FromBytes(e.CommunityID.Bytes[:])
		result["community_id"] = cID.String()
	}
	if e.EndTime.Valid {
		result["end_time"] = e.EndTime.Time
	}
	if e.MatchFormat.Valid {
		result["match_format"] = string(e.MatchFormat.MatchFormat)
	}

	return result
}

func buildEventListItem(e repository.ListEventsRow) map[string]interface{} {
	eID, _ := uuid.FromBytes(e.ID.Bytes[:])
	creatorID, _ := uuid.FromBytes(e.CreatedBy.Bytes[:])

	result := map[string]interface{}{
		"id":                   eID.String(),
		"title":                e.Title,
		"event_type":           string(e.EventType),
		"status":               string(e.Status.EventStatus),
		"start_time":           e.StartTime.Time,
		"max_participants":     e.MaxParticipants.Int32,
		"current_participants": e.CurrentParticipants.Int32,
		"min_level":            numericToFloat(e.MinLevel),
		"max_level":            numericToFloat(e.MaxLevel),
		"location_name":        e.LocationName.String,
		"is_paid":              e.IsPaid.Bool,
		"created_by":           creatorID.String(),
	}

	if e.CommunityID.Valid {
		cID, _ := uuid.FromBytes(e.CommunityID.Bytes[:])
		result["community_id"] = cID.String()
	}

	return result
}
