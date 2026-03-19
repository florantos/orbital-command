package domain_test

import (
	"testing"

	"github.com/florantos/orbital-command/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewAuditEvent_ValidAuditEventCreatesSuccessfully(t *testing.T) {
	av := domain.NewAuditEvent("module.registered", "module", "abc-123", "Commander Chen", "Registered module: Navigation Array")

	assert.Equal(t, "module.registered", av.Action)
	assert.Equal(t, "module", av.EntityType)
	assert.Equal(t, "abc-123", av.EntityId)
	assert.Equal(t, "Commander Chen", av.Actor)
	assert.Equal(t, "Registered module: Navigation Array", av.Detail)

}
