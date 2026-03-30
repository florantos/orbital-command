package domain_test

import (
	"testing"

	"github.com/florantos/orbital-command/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCrewMember_ValidCrewMemberCreatesSuccessfully(t *testing.T) {
	quals := []domain.Capability{domain.CapabilityDocking, domain.CapabilityAtmosphereRecycling}
	cm, err := domain.NewCrewMember("test crew member", domain.RoleEngineer, quals)
	require.NoError(t, err)

	assert.Equal(t, "test crew member", cm.Name)
	assert.Equal(t, domain.RoleEngineer, cm.Role)
	assert.Equal(t, quals, cm.Qualifications)

}

func TestNewCrewMember_NameValidation(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr string
	}{
		{
			name:        "empty name returns error",
			input:       "",
			expectedErr: "name is required",
		},
		{
			name:        "whitespace only returns error",
			input:       " ",
			expectedErr: "name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewCrewMember(tt.input, domain.RoleEngineer, []domain.Capability{domain.CapabilityDocking})

			require.Error(t, err)
			var ve *domain.ValidationError
			require.ErrorAs(t, err, &ve)
			assert.Equal(t, tt.expectedErr, ve.Fields["name"])
		})
	}
}

func TestNewCrewMember_RoleValidation(t *testing.T) {
	tests := []struct {
		name        string
		input       domain.Role
		expectedErr string
	}{
		{
			name:        "invalid role returns error",
			input:       "wood-chopper",
			expectedErr: "wood-chopper is not a valid role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewCrewMember("jon chop", tt.input, []domain.Capability{domain.CapabilityDocking})

			require.Error(t, err)
			var ve *domain.ValidationError
			require.ErrorAs(t, err, &ve)
			assert.Equal(t, tt.expectedErr, ve.Fields["role"])
		})
	}
}

func TestNewCrewMember_QualificationsValidation(t *testing.T) {
	tests := []struct {
		name        string
		input       []domain.Capability
		expectedErr string
	}{
		{
			name:        "0 qualifications returns error",
			input:       []domain.Capability{},
			expectedErr: "at least one qualification is required",
		},
		{
			name:        "invalid qualification returns error",
			input:       []domain.Capability{"wood-chopper"},
			expectedErr: "wood-chopper is not a valid qualification",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewCrewMember("jon chop", domain.RoleEngineer, tt.input)

			require.Error(t, err)
			var ve *domain.ValidationError
			require.ErrorAs(t, err, &ve)
			assert.Equal(t, tt.expectedErr, ve.Fields["qualifications"])
		})
	}
}
