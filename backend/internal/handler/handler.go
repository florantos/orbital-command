package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Handler struct {
	logger     *slog.Logger
	moduleRepo ModuleRepository
}

func NewHandler(logger *slog.Logger, moduleRepo ModuleRepository) *Handler {
	return &Handler{
		logger:     logger,
		moduleRepo: moduleRepo,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	body, _ := json.Marshal(ErrorResponse{Error: message})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body) //nolint:errcheck
}
