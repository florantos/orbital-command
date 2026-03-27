package repository_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	testPool, err = pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to create test pool: %v", err)
	}

	code := m.Run()

	testPool.Close()
	os.Exit(code)
}
