package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/florantos/orbital-command/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

type ModuleRepo struct {
	db DBTX
}

func NewModuleRepo(db DBTX) *ModuleRepo {
	return &ModuleRepo{db: db}
}

func (r *ModuleRepo) Create(ctx context.Context, module *domain.Module) (*domain.Module, error) {
	query := `
		INSERT INTO modules (name, description, health_state)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, health_state, created_at, updated_at
	`

	created := &domain.Module{}

	err := r.db.QueryRow(ctx, query, module.Name, module.Description, module.HealthState).Scan(
		&created.ID,
		&created.Name,
		&created.Description,
		&created.HealthState,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("create module: %w", domain.ErrDuplicateModuleName)
		}
		return nil, fmt.Errorf("create module: %w", err)
	}

	return created, nil
}
