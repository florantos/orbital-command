package repository_test

import (
	"context"
	"testing"

	"github.com/florantos/orbital-command/internal/domain"
	"github.com/florantos/orbital-command/internal/repository"
	"github.com/florantos/orbital-command/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditEventRepo_Create_PersistsAuditEvent(t *testing.T) {
	pool := testutil.NewTestPool(t)
	tx := testutil.NewTestTx(t, pool)
	repo := repository.NewAuditEventRepo(tx)

	event := domain.NewAuditEvent("module.registered", "module", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "Commander Chen", "Registered module: Navigation Array")

	err := repo.Create(context.Background(), event)

	require.NoError(t, err)

}

func TestAuditEventRepo_ReadAll_ReturnsAllEvents(t *testing.T) {
	pool := testutil.NewTestPool(t)
	tx := testutil.NewTestTx(t, pool)
	repo := repository.NewAuditEventRepo(tx)

	event := domain.NewAuditEvent("module.registered", "module", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12", "System", "Registered module: Navigation Array2")
	event2 := domain.NewAuditEvent("module.registered", "module", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "Commander Chen", "Registered module: Navigation Array")

	err := repo.Create(t.Context(), event)
	require.NoError(t, err)
	err = repo.Create(t.Context(), event2)
	require.NoError(t, err)

	events, err := repo.ReadAll(context.Background())
	require.NoError(t, err)

	assert.Len(t, events, 2)

}

func TestAuditEventRepo_ReadAll_ReturnsEmptyArrayWhenNoEvents(t *testing.T) {
	pool := testutil.NewTestPool(t)
	tx := testutil.NewTestTx(t, pool)
	repo := repository.NewAuditEventRepo(tx)

	events, err := repo.ReadAll(context.Background())
	require.NoError(t, err)

	assert.NotNil(t, events)
	assert.Len(t, events, 0)
}
