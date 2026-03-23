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
		return fmt.Errorf("create audit event: %w", err)
	}

	return nil
}

func (r *AuditEventRepo) ReadAll(ctx context.Context) ([]domain.AuditEvent, error) {
	query := `
		SELECT id, action, entity_type, entity_id, actor, detail, occurred_at
		FROM audit_events
		ORDER BY occurred_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("read all audit events: %w", err)
	}
	defer rows.Close()

	events := []domain.AuditEvent{}
	for rows.Next() {
		var e domain.AuditEvent
		err := rows.Scan(&e.ID, &e.Action, &e.EntityType, &e.EntityID, &e.Actor, &e.Detail, &e.OccurredAt)
		if err != nil {
			return nil, fmt.Errorf("read all audit events: scan: %w", err)
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read all audit events: rows: %w", err)
	}

	return events, nil
}
