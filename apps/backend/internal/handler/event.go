package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/service"
)

// EventHandler handles event endpoints
type EventHandler struct {
	eventService *service.EventService
}

// NewEventHandler creates a new EventHandler
func NewEventHandler(eventService *service.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

// Create handles POST /v1/events
func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	var input service.CreateEventInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if input.Title == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "title is required")
		return
	}
	if input.EventType == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "event_type is required")
		return
	}
	if input.StartTime.IsZero() {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "start_time is required")
		return
	}

	event, err := h.eventService.Create(r.Context(), userID, input)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, event)
}

// List handles GET /v1/events
func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	input := service.ListEventsInput{
		EventType:   q.Get("event_type"),
		Status:      q.Get("status"),
		Composition: q.Get("composition"),
		CommunityID: q.Get("community_id"),
		District:    q.Get("district"),
		Sort:        q.Get("sort"),
		Page:        queryInt(q.Get("page"), 1),
		PerPage:     queryInt(q.Get("per_page"), 20),
		MinLevel:    parseQueryFloat(q.Get("min_level")),
		MaxLevel:    parseQueryFloat(q.Get("max_level")),
	}

	if v := q.Get("date_from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			input.DateFrom = &t
		}
	}
	if v := q.Get("date_to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			input.DateTo = &t
		}
	}

	events, pagination, err := h.eventService.List(r.Context(), input)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondPaginated(w, http.StatusOK, events, *pagination)
}

// GetByID handles GET /v1/events/:id
func (h *EventHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	eventID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID")
		return
	}

	event, err := h.eventService.GetByID(r.Context(), userID, eventID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, event)
}

// Join handles POST /v1/events/:id/join
func (h *EventHandler) Join(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	eventID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID")
		return
	}

	result, err := h.eventService.Join(r.Context(), userID, eventID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// Leave handles POST /v1/events/:id/leave
func (h *EventHandler) Leave(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	eventID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID")
		return
	}

	if err := h.eventService.Leave(r.Context(), userID, eventID); err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Successfully left event"})
}

// UpdateStatus handles PATCH /v1/events/:id/status
func (h *EventHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	eventID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID")
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if body.Status == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "status is required")
		return
	}

	event, err := h.eventService.UpdateStatus(r.Context(), userID, eventID, body.Status)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, event)
}

// GetCalendar handles GET /v1/events/calendar
func (h *EventHandler) GetCalendar(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	year, err := strconv.Atoi(q.Get("year"))
	if err != nil || year < 2020 || year > 2100 {
		year = time.Now().Year()
	}

	month, err := strconv.Atoi(q.Get("month"))
	if err != nil || month < 1 || month > 12 {
		month = int(time.Now().Month())
	}

	communityID := q.Get("community_id")

	calendar, err := h.eventService.GetCalendar(r.Context(), year, month, communityID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, calendar)
}

// GetMyEvents handles GET /v1/events/my
func (h *EventHandler) GetMyEvents(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	q := r.URL.Query()
	tab := q.Get("tab")
	if tab == "" {
		tab = "joined"
	}
	page := queryInt(q.Get("page"), 1)
	perPage := queryInt(q.Get("per_page"), 20)

	events, err := h.eventService.GetMyEvents(r.Context(), userID, tab, page, perPage)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, events)
}

// ListParticipants handles GET /v1/events/:id/participants
func (h *EventHandler) ListParticipants(w http.ResponseWriter, r *http.Request) {
	eventID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID")
		return
	}

	// Reuse GetByID which includes participants
	userID, _ := getUserUUID(r)
	event, err := h.eventService.GetByID(r.Context(), userID, eventID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, event["participants"])
}
