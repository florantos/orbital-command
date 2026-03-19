package repository

import (
	"context"
	"fmt"

	"github.com/florantos/orbital-command/internal/domain"
)

type AuditEventRepo struct {
	db DBTX
}

func NewAuditEventRepo(db DBTX) *AuditEventRepo {
	return &AuditEventRepo{db: db}

}

func (r *AuditEventRepo) Create(ctx context.Context, event *domain.AuditEvent) error {
	query := `
		INSERT INTO audit_events (action, entity_type, entity_id, actor, detail) 
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query, event.Action, event.EntityType, event.EntityID, event.Actor, event.Detail)
	if err != nil {
		return fmt.Errorf("create audit_event: %w", err)
	}

	return nil
}
