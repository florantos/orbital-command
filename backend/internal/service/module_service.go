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

type ModuleRepository interface {
	Create(ctx context.Context, db database.DBTX, cm *domain.Module) (*domain.Module, error)
	ReadAll(ctx context.Context, db database.DBTX) ([]domain.Module, error)
}

type ModuleService struct {
	pool       *pgxpool.Pool
	logger     *slog.Logger
	moduleRepo ModuleRepository
	auditRepo  AuditEventRepository
}

func NewModuleService(pool *pgxpool.Pool, logger *slog.Logger, moduleRepo ModuleRepository, auditRepo AuditEventRepository) *ModuleService {
	return &ModuleService{pool: pool, logger: logger, moduleRepo: moduleRepo, auditRepo: auditRepo}
}

func (s *ModuleService) Create(ctx context.Context, name string, description string) (*domain.Module, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("create module: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			s.logger.Error("failed to rollback transaction", "error", err)

		}
	}()

	module, err := domain.NewModule(name, description)
	if err != nil {
		return nil, fmt.Errorf("create module: %w", err)
	}

	created, err := s.moduleRepo.Create(ctx, tx, module)
	if err != nil {
		return nil, fmt.Errorf("create module: %w", err)
	}

	event := domain.NewAuditEvent("module.registered", "module", created.ID, "Commander", fmt.Sprintf("Registered module: %s", created.Name))
	if err := s.auditRepo.Create(ctx, tx, event); err != nil {
		return nil, fmt.Errorf("failed to create audit event: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("create module: %w", err)
	}
	return created, nil
}

func (s *ModuleService) ReadAll(ctx context.Context) ([]domain.Module, error) {
	modules, err := s.moduleRepo.ReadAll(ctx, s.pool)
	if err != nil {
		return nil, fmt.Errorf("read all module: %w", err)
	}
	return modules, nil
}
