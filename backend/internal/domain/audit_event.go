package domain

import (
	"time"
)

type AuditEvent struct {
	ID         string
	Action     string
	EntityType string
	EntityID   string
	Actor      string
	Detail     string
	OccurredAt time.Time
}

func NewAuditEvent(action, entityType, entityId, actor, detail string) *AuditEvent {

	return &AuditEvent{
		Action:     action,
		EntityType: entityType,
		EntityID:   entityId,
		Actor:      actor,
		Detail:     detail,
	}
}
