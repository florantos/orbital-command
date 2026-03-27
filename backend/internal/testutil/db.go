package testutil

import (
	"context"
	"testing"

	"github.com/florantos/orbital-command/internal/database"
	"github.com/florantos/orbital-command/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func NewTestTx(t *testing.T, pool *pgxpool.Pool) pgx.Tx {
	t.Helper()

	tx, err := pool.Begin(context.Background())
	require.NoError(t, err)

	t.Cleanup(func() {
		err := tx.Rollback(context.Background())
		if err != nil && err != pgx.ErrTxClosed {
			t.Errorf("failed to rollback test transaction: %v", err)
		}
	})

	return tx
}

func NewTestModule(t *testing.T, opts ...func(*domain.Module)) *domain.Module {
	t.Helper()
	m, err := domain.NewModule(
		"Test Module"+uuid.NewString()[:8],
		"Test Description"+uuid.NewString()[:8],
	)
	require.NoError(t, err)
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func SeedModules(t *testing.T, db database.DBTX, modules []*domain.Module) {
	query := `
		INSERT INTO modules (name, description, health_state) 
		VALUES ($1, $2, $3)
	`
	for _, m := range modules {
		_, err := db.Exec(context.Background(), query, m.Name, m.Description, m.HealthState)
		if err != nil {
			t.Fatalf("seeding modules: %v", err)
		}
	}
}

func NewTestAuditEvent(t *testing.T, opts ...func(*domain.AuditEvent)) *domain.AuditEvent {
	t.Helper()
	ae := domain.NewAuditEvent(
		"module.registered",
		"module",
		uuid.NewString(),
		"Commander Chen",
		"Registered module: Navigation Array",
	)
	for _, opt := range opts {
		opt(ae)
	}
	return ae
}

func SeedAuditEvents(t *testing.T, db database.DBTX, auditEvents []*domain.AuditEvent) {
	t.Helper()
	query := `
		INSERT INTO audit_events (action, entity_type, entity_id, actor, detail) 
		VALUES ($1, $2, $3, $4, $5)
	`
	for _, ae := range auditEvents {
		_, err := db.Exec(context.Background(), query, ae.Action, ae.EntityType, ae.EntityID, ae.Actor, ae.Detail)
		if err != nil {
			t.Fatalf("seeding audit events: %v", err)
		}
	}
}
