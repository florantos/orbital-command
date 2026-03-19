package domain

import (
	"time"
)

type AuditEvent struct {
	ID         string
	Action     string
	EntityType string
	EntityId   string
	Actor      string
	Detail     string
	OccuredAt  time.Time
}

func NewAuditEvent(action, entityType, entityId, actor, detail string) *AuditEvent {

	return &AuditEvent{
		Action:     action,
		EntityType: entityType,
		EntityId:   entityId,
		Actor:      actor,
		Detail:     detail,
	}
}
