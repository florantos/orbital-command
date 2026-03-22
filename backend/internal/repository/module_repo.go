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

func (r *ModuleRepo) ReadAll(ctx context.Context) ([]domain.Module, error) {
	query := `
		SELECT id, name, description, health_state, created_at
		FROM modules
		ORDER BY
			Case health_state
				WHEN 'offline'	    THEN 1
				WHEN 'unresponsive' THEN 2
				WHEN 'critical'		THEN 3
				WHEN 'degraded'		THEN 4
				WHEN 'operation'	THEN 5
		END
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("read all modules: %w", err)
	}
	defer rows.Close()

	modules := []domain.Module{}
	for rows.Next() {
		var m domain.Module
		err := rows.Scan(&m.ID, &m.Name, &m.Description, &m.HealthState, &m.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("read all modules: scan: %w", err)
		}
		modules = append(modules, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read all modules: rows: %w", err)
	}

	return modules, nil
}
