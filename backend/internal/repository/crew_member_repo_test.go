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

func TestCrewRepo_ReadAll_ReturnsAllCrewMembers(t *testing.T) {
	tx := testutil.NewTestTx(t, testPool)
	repo := repository.NewCrewRepo()

	crew := make([]*domain.CrewMember, 10)
	for i := range crew {
		crew[i] = testutil.NewTestCrewMember(t)
	}
	testutil.SeedCrewMembers(t, tx, crew)

	result, err := repo.ReadAll(context.Background(), tx)
	require.NoError(t, err)

	assert.Len(t, result, 10)

}
func TestCrewRepo_ReadAll_ReturnsEmptyArrayWhenNoCrewMembers(t *testing.T) {
	repo := repository.NewCrewRepo()

	crew, err := repo.ReadAll(context.Background(), testPool)
	require.NoError(t, err)

	assert.NotNil(t, crew)
	assert.Len(t, crew, 0)
}
