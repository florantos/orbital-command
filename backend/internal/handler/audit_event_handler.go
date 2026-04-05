package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/florantos/orbital-command/internal/domain"
)

type AuditEventHandler struct {
	logger            *slog.Logger
	auditEventService AuditEventService
}

func NewAuditHandler(logger *slog.Logger, auditEventService AuditEventService) *AuditEventHandler {
	return &AuditEventHandler{
		logger:            logger,
		auditEventService: auditEventService,
	}
}

type AuditEventService interface {
	ReadAll(ctx context.Context) ([]domain.AuditEvent, error)
}
type AuditEventResponse struct {
	ID         string    `json:"id"`
	Action     string    `json:"action"`
	EntityType string    `json:"entityType"`
	EntityID   string    `json:"entityID"`
	Actor      string    `json:"actor"`
	Detail     string    `json:"detail"`
	OccurredAt time.Time `json:"occurredAt"`
}
type ReadAllAuditEventsResponse struct {
	AuditEvents []AuditEventResponse `json:"auditEvents"`
}

func (h *AuditEventHandler) ReadAllAuditEvents(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("reading all audit events")

	events, err := h.auditEventService.ReadAll(r.Context())
	if err != nil {
		h.logger.Error("failed to read all audit events", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := ReadAllAuditEventsResponse{
		AuditEvents: []AuditEventResponse{},
	}

	for _, e := range events {
		resp.AuditEvents = append(resp.AuditEvents, AuditEventResponse{
			ID:         e.ID,
			Action:     e.Action,
			EntityType: e.EntityType,
			EntityID:   e.EntityID,
			Actor:      e.Actor,
			Detail:     e.Detail,
			OccurredAt: e.OccurredAt,
		})
	}
	h.logger.Info("audit events read", "count", len(events))

	err = writeJSON(w, http.StatusOK, resp)
	if err != nil {
		h.logger.Error("failed to marshal response", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
}
