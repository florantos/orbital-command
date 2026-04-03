package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/florantos/orbital-command/internal/database"
	"github.com/florantos/orbital-command/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditEventRepository interface {
	Create(ctx context.Context, db database.DBTX, event *domain.AuditEvent) error
	ReadAll(ctx context.Context, db database.DBTX) ([]domain.AuditEvent, error)
}

type Handler struct {
	logger         *slog.Logger
	pool           pgxpool.Pool
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
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
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

func writeValidationError(w http.ResponseWriter, ve *domain.ValidationError) {
	body, err := json.Marshal(ErrorResponse{Error: "validation failed", Fields: ve.Fields})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
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
