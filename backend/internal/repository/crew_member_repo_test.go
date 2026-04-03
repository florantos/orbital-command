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

func TestCrewRepo_Create_PersistsAndReturnsCrewMember(t *testing.T) {
	tx := testutil.NewTestTx(t, testPool)
	repo := repository.NewCrewRepo()

	cm := testutil.NewTestCrewMember(t)

	created, err := repo.Create(context.Background(), tx, cm)
	require.NoError(t, err)

	assert.NotEmpty(t, created.ID)
	assert.Equal(t, cm.Name, created.Name)
	assert.Equal(t, cm.Role, created.Role)
	assert.Equal(t, cm.Qualifications, created.Qualifications)
	assert.NotZero(t, created.CreatedAt)
	assert.NotZero(t, created.UpdatedAt)
}

func TestCrewRepo_Create_ReturnsErrorOnDuplicateName(t *testing.T) {
	tx := testutil.NewTestTx(t, testPool)
	repo := repository.NewCrewRepo()

	cm := testutil.NewTestCrewMember(t)

	_, err := repo.Create(context.Background(), tx, cm)
	require.NoError(t, err)

	_, err = repo.Create(context.Background(), tx, cm)
	assert.ErrorIs(t, err, domain.ErrDuplicateCrewMemberName)
}
