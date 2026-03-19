package repository_test

import (
	"context"
	"testing"

	"github.com/florantos/orbital-command/internal/domain"
	"github.com/florantos/orbital-command/internal/repository"
	"github.com/florantos/orbital-command/internal/testutil"
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
