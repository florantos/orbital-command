package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/florantos/orbital-command/internal/database"
	"github.com/florantos/orbital-command/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

type CrewRepo struct{}

func NewCrewRepo() *CrewRepo {
	return &CrewRepo{}
}

func (r *CrewRepo) Create(ctx context.Context, db database.DBTX, cm *domain.CrewMember) (*domain.CrewMember, error) {
	query := `
		INSERT INTO crew (name, role)
		VALUES ($1, $2)
		RETURNING id, name, role, created_at, updated_at
	`

	created := &domain.CrewMember{}

	err := db.QueryRow(ctx, query, cm.Name, cm.Role).Scan(
		&created.ID,
		&created.Name,
		&created.Role,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("create crew member: %w", domain.ErrDuplicateCrewMemberName)
		}
		return nil, fmt.Errorf("created crew member: %w", err)
	}

	capNames := make([]string, len(cm.Qualifications))
	for i, cap := range cm.Qualifications {
		capNames[i] = string(cap)
	}

	query = `
		INSERT INTO crew_capabilities (crew_id, name)
    	SELECT $1, unnest($2::text[])
	`
	_, err = db.Exec(ctx, query, created.ID, capNames)

	if err != nil {
		return nil, fmt.Errorf("create crew member capabilities: %w", err)
	}
	created.Qualifications = cm.Qualifications
	return created, nil
}
