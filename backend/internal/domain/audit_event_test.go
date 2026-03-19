package domain_test

import (
	"testing"

	"github.com/florantos/orbital-command/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewAuditEvent_ValidAuditEventCreatesSuccessfully(t *testing.T) {
	av := domain.NewAuditEvent("module.registered", "module", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "Commander Chen", "Registered module: Navigation Array")

	assert.Equal(t, "module.registered", av.Action)
	assert.Equal(t, "module", av.EntityType)
	assert.Equal(t, "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", av.EntityID)
	assert.Equal(t, "Commander Chen", av.Actor)
	assert.Equal(t, "Registered module: Navigation Array", av.Detail)

}
