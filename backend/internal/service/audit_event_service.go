package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/florantos/orbital-command/internal/database"
	"github.com/florantos/orbital-command/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditEventRepository interface {
	Create(ctx context.Context, db database.DBTX, event *domain.AuditEvent) error
	ReadAll(ctx context.Context, db database.DBTX) ([]domain.AuditEvent, error)
}

type AuditEventService struct {
	pool      *pgxpool.Pool
	logger    *slog.Logger
	auditRepo AuditEventRepository
}

func NewAuditEventService(pool *pgxpool.Pool, logger *slog.Logger, auditRepo AuditEventRepository) *AuditEventService {
	return &AuditEventService{pool: pool, logger: logger, auditRepo: auditRepo}
}

func (s *AuditEventService) ReadAll(ctx context.Context) ([]domain.AuditEvent, error) {
	events, err := s.auditRepo.ReadAll(ctx, s.pool)
	if err != nil {
		return nil, fmt.Errorf("read all audit events: %w", err)
	}
	return events, nil
}
