package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/florantos/orbital-command/internal/domain"
)

type AuditEventRepository interface {
	Create(ctx context.Context, event *domain.AuditEvent) error
	ReadAll(ctx context.Context) ([]domain.AuditEvent, error)
}

type Handler struct {
	logger         *slog.Logger
	moduleRepo     ModuleRepository
	auditEventRepo AuditEventRepository
}

func NewHandler(logger *slog.Logger, moduleRepo ModuleRepository, auditEventRepo AuditEventRepository) *Handler {
	return &Handler{
		logger:         logger,
		moduleRepo:     moduleRepo,
		auditEventRepo: auditEventRepo,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	body, err := json.Marshal(ErrorResponse{Error: message})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
	return nil
}
