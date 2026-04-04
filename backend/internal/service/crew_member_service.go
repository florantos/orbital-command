package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/florantos/orbital-command/internal/database"
	"github.com/florantos/orbital-command/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CrewRepository interface {
	Create(ctx context.Context, db database.DBTX, cm *domain.CrewMember) (*domain.CrewMember, error)
}

type CrewService struct {
	pool      *pgxpool.Pool
	logger    *slog.Logger
	crewRepo  CrewRepository
	auditRepo AuditEventRepository
}

func NewCrewService(pool *pgxpool.Pool, logger *slog.Logger, crewRepo CrewRepository, auditRepo AuditEventRepository) *CrewService {
	return &CrewService{pool: pool, logger: logger, crewRepo: crewRepo, auditRepo: auditRepo}
}

func (s *CrewService) Create(ctx context.Context, name string, role domain.Role, qualifications []domain.Capability) (*domain.CrewMember, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("create crew member: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			s.logger.Error("failed to rollback transaction", "error", err)

		}
	}()
	s.logger.Info("creating crew member", "name", name)

	crewMember, err := domain.NewCrewMember(name, role, qualifications)
	if err != nil {
		return nil, fmt.Errorf("create crew member: %w", err)
	}

	created, err := s.crewRepo.Create(ctx, tx, crewMember)
	if err != nil {
		return nil, fmt.Errorf("create crew member: %w", err)
	}

	event := domain.NewAuditEvent("crew.registered", "crew", created.ID, "Commander", fmt.Sprintf("Registered crew member: %s", created.Name))
	if err := s.auditRepo.Create(ctx, tx, event); err != nil {
		return nil, fmt.Errorf("failed to create audit event: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("create crew member: %w", err)
	}
	return created, nil
}
