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
	tx := testutil.NewTestTx(t, testPool)
	repo := repository.NewAuditEventRepo(tx)

	event := testutil.NewTestAuditEvent(t)
	err := repo.Create(context.Background(), event)

	require.NoError(t, err)

}

func TestAuditEventRepo_ReadAll_ReturnsAllEvents(t *testing.T) {
	tx := testutil.NewTestTx(t, testPool)
	repo := repository.NewAuditEventRepo(tx)

	events := make([]*domain.AuditEvent, 10)
	for i := range events {
		events[i] = testutil.NewTestAuditEvent(t)
	}
	testutil.SeedAuditEvents(t, tx, events)

	result, err := repo.ReadAll(context.Background())
	require.NoError(t, err)
	assert.Len(t, result, 10)

}

func TestAuditEventRepo_ReadAll_ReturnsEmptyArrayWhenNoEvents(t *testing.T) {
	tx := testutil.NewTestTx(t, testPool)
	repo := repository.NewAuditEventRepo(tx)

	events, err := repo.ReadAll(context.Background())
	require.NoError(t, err)

	assert.NotNil(t, events)
	assert.Len(t, events, 0)
}
