package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/florantos/orbital-command/internal/database"
	"github.com/florantos/orbital-command/internal/domain"
	"github.com/florantos/orbital-command/internal/repository"
	"github.com/florantos/orbital-command/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockAuditEventRepo struct {
	createFn func(ctx context.Context, db database.DBTX, event *domain.AuditEvent) error
}

func (m *mockAuditEventRepo) Create(ctx context.Context, db database.DBTX, event *domain.AuditEvent) error {
	return m.createFn(ctx, db, event)
}

func TestCrewService_Create_Success(t *testing.T) {
	crewRepo := repository.NewCrewRepo()
	auditRepo := repository.NewAuditEventRepo()

	svc := service.NewCrewService(testPool, crewRepo, auditRepo)

	name := "John Snow " + uuid.NewString()[:8]

	created, err := svc.Create(
		context.Background(),
		name,
		domain.RoleEngineer,
		[]domain.Capability{domain.CapabilityDocking, domain.CapabilityNavigation},
	)

	require.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, name, created.Name)
	assert.Equal(t, domain.RoleEngineer, created.Role)
	assert.Len(t, created.Qualifications, 2)

	t.Cleanup(func() {
		_, err := testPool.Exec(context.Background(), "DELETE FROM crew WHERE id = $1", created.ID)
		if err != nil {
			t.Errorf("failed to clean up crew member: %v", err)
		}
	})

}

func TestCrewService_Create_RollsBackOnAuditFailure(t *testing.T) {
	crewRepo := repository.NewCrewRepo()

	auditRepo := &mockAuditEventRepo{
		createFn: func(ctx context.Context, db database.DBTX, event *domain.AuditEvent) error {
			return fmt.Errorf("audit db error")
		},
	}

	svc := service.NewCrewService(testPool, crewRepo, auditRepo)

	name := "John Snow " + uuid.NewString()[:8]
	_, err := svc.Create(
		context.Background(),
		name,
		domain.RoleEngineer,
		[]domain.Capability{domain.CapabilityDocking, domain.CapabilityNavigation},
	)

	require.Error(t, err)

	var count int
	err = testPool.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM crew WHERE name = $1",
		name,
	).Scan(&count)

	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestCrewService_Create_RollsBackOnCrewFailure(t *testing.T) {
	crewRepo := repository.NewCrewRepo()
	auditRepo := repository.NewAuditEventRepo()

	svc := service.NewCrewService(testPool, crewRepo, auditRepo)

	name := "John Snow " + uuid.NewString()[:8]

	firstCreated, err := svc.Create(
		context.Background(),
		name,
		domain.RoleEngineer,
		[]domain.Capability{domain.CapabilityDocking},
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err := testPool.Exec(context.Background(), "DELETE FROM crew WHERE id = $1", firstCreated.ID)
		if err != nil {
			t.Errorf("failed to clean up crew member: %v", err)
		}
	})

	_, err = svc.Create(
		context.Background(),
		name,
		domain.RoleEngineer,
		[]domain.Capability{domain.CapabilityDocking},
	)
	require.Error(t, err)

	var count int
	err = testPool.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM crew WHERE name = $1",
		name,
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	var auditCount int
	err = testPool.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM audit_events WHERE detail LIKE $1",
		"%"+name+"%",
	).Scan(&auditCount)
	require.NoError(t, err)
	assert.Equal(t, 1, auditCount)
}
