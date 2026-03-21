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
	body, _ := json.Marshal(ErrorResponse{Error: message})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}
