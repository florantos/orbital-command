package testutil

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func NewTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	pool, err := pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

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
