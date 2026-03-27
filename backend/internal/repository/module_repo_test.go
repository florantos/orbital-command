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

func TestModuleRepo_Create_PersistsAndReturnsModule(t *testing.T) {
	tx := testutil.NewTestTx(t, testPool)
	repo := repository.NewModuleRepo(tx)

	module := testutil.NewTestModule(t)

	created, err := repo.Create(context.Background(), module)

	require.NoError(t, err)

	assert.NotEmpty(t, created.ID)
	assert.Equal(t, module.Name, created.Name)
	assert.Equal(t, module.Description, created.Description)
	assert.Equal(t, module.HealthState, created.HealthState)
	assert.NotZero(t, created.CreatedAt)
	assert.NotZero(t, created.UpdatedAt)
}

func TestModuleRepo_Create_ReturnsErrorOnDuplicateName(t *testing.T) {
	tx := testutil.NewTestTx(t, testPool)
	repo := repository.NewModuleRepo(tx)

	module := testutil.NewTestModule(t)

	_, err := repo.Create(context.Background(), module)
	require.NoError(t, err)

	_, err = repo.Create(context.Background(), module)
	assert.ErrorIs(t, err, domain.ErrDuplicateModuleName)

}

func TestModuleRepo_ReadAll_ReturnsAllModules(t *testing.T) {
	tx := testutil.NewTestTx(t, testPool)
	repo := repository.NewModuleRepo(tx)

	modules := make([]*domain.Module, 10)
	for i := range modules {
		modules[i] = testutil.NewTestModule(t)
	}
	testutil.SeedModules(t, tx, modules)

	result, err := repo.ReadAll(context.Background())
	require.NoError(t, err)

	assert.Len(t, result, 10)

}
func TestModuleRepo_ReadAll_ReturnsEmptyArrayWhenNoModules(t *testing.T) {
	tx := testutil.NewTestTx(t, testPool)
	repo := repository.NewModuleRepo(tx)

	modules, err := repo.ReadAll(context.Background())
	require.NoError(t, err)

	assert.NotNil(t, modules)
	assert.Len(t, modules, 0)
}
