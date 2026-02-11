package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/pkg/validator"
	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/service"
)

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{"data": data})
}

func respondError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	})
}

func respondValidationError(w http.ResponseWriter, err error) {
	details := validator.FormatErrors(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"code":    "VALIDATION_ERROR",
			"message": "Validation failed",
			"details": details,
		},
	})
}

func handleServiceError(w http.ResponseWriter, err error) {
	if appErr, ok := service.IsAppError(err); ok {
		respondError(w, appErr.Status, appErr.Code, appErr.Error())
		return
	}

	// Unhandled error - log and return 500
	slog.Error("unhandled service error", "error", err)
	respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
}
