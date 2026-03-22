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
	pool := testutil.NewTestPool(t)
	tx := testutil.NewTestTx(t, pool)
	repo := repository.NewModuleRepo(tx)

	module, err := domain.NewModule("Navigation Array", "Controls navigation systems")
	require.NoError(t, err)

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
	pool := testutil.NewTestPool(t)
	tx := testutil.NewTestTx(t, pool)
	repo := repository.NewModuleRepo(tx)

	module, err := domain.NewModule("Navigation Array", "Controls navigation systems")
	require.NoError(t, err)

	_, err = repo.Create(context.Background(), module)
	require.NoError(t, err)

	_, err = repo.Create(context.Background(), module)
	assert.ErrorIs(t, err, domain.ErrDuplicateModuleName)

}

func TestModuleRepo_ReadAll_ReturnsAllModules(t *testing.T) {
	pool := testutil.NewTestPool(t)
	tx := testutil.NewTestTx(t, pool)
	repo := repository.NewModuleRepo(tx)

	module, err := domain.NewModule("Navigation Array", "Controls navigation systems")
	require.NoError(t, err)
	_, err = repo.Create(context.Background(), module)
	require.NoError(t, err)

	module2, err := domain.NewModule("Navigation Array2", "Controls navigation systems2")
	require.NoError(t, err)
	_, err = repo.Create(context.Background(), module2)
	require.NoError(t, err)

	modules, err := repo.ReadAll(context.Background())
	require.NoError(t, err)

	assert.Len(t, modules, 2)

}
func TestModuleRepo_ReadAll_ReturnsEmptyArrayWhenNoModules(t *testing.T) {
	pool := testutil.NewTestPool(t)
	tx := testutil.NewTestTx(t, pool)
	repo := repository.NewModuleRepo(tx)

	modules, err := repo.ReadAll(context.Background())
	require.NoError(t, err)

	assert.NotNil(t, modules)
	assert.Len(t, modules, 0)
}
