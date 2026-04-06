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
		return nil, fmt.Errorf("create crew member: %w", err)
	}

	capNames := make([]string, len(cm.Qualifications))
	for i, cap := range cm.Qualifications {
		capNames[i] = string(cap)
	}

	query = `
		INSERT INTO crew_qualifications (crew_id, name)
    	SELECT $1, unnest($2::text[])
	`
	_, err = db.Exec(ctx, query, created.ID, capNames)

	if err != nil {
		return nil, fmt.Errorf("create crew member capabilities: %w", err)
	}
	created.Qualifications = cm.Qualifications
	return created, nil
}

func (r *CrewRepo) ReadAll(ctx context.Context, db database.DBTX) ([]domain.CrewMember, error) {
	query := `
        SELECT id, name, role, created_at
        FROM crew
    `
	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("read all crew members: %w", err)
	}
	defer rows.Close()

	crew := []domain.CrewMember{}
	for rows.Next() {
		var cm domain.CrewMember
		err := rows.Scan(&cm.ID, &cm.Name, &cm.Role, &cm.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("read all crew members: scan: %w", err)
		}
		crew = append(crew, cm)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read all crew members: rows: %w", err)
	}
	rows.Close()

	capQuery := `
        SELECT crew_id, name
        FROM crew_qualifications
    `
	capRows, err := db.Query(ctx, capQuery)
	if err != nil {
		return nil, fmt.Errorf("read all crew capabilities: %w", err)
	}
	defer capRows.Close()

	capMap := make(map[string][]domain.Capability)
	for capRows.Next() {
		var crewID string
		var cap domain.Capability
		err := capRows.Scan(&crewID, &cap)
		if err != nil {
			return nil, fmt.Errorf("read all crew capabilities: scan: %w", err)
		}
		capMap[crewID] = append(capMap[crewID], cap)
	}
	if err := capRows.Err(); err != nil {
		return nil, fmt.Errorf("read all crew capabilities: rows: %w", err)
	}

	for i, cm := range crew {
		crew[i].Qualifications = capMap[cm.ID]
	}

	return crew, nil
}
