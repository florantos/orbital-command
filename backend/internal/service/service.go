package service

import (
	"context"

	"github.com/florantos/orbital-command/internal/database"
	"github.com/florantos/orbital-command/internal/domain"
)

type AuditEventRepository interface {
	Create(ctx context.Context, db database.DBTX, event *domain.AuditEvent) error
}
