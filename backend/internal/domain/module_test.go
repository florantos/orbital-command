package domain_test

import (
	"testing"

	"github.com/florantos/orbital-command/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewModule_ValidModuleCreatesSuccessfully(t *testing.T) {

	module, err := domain.NewModule("test module", "test description")

	require.NoError(t, err)

	assert.Equal(t, "test module", module.Name)
	assert.Equal(t, "test description", module.Description)
	assert.Equal(t, domain.HealthStateOperational, module.HealthState)
}

func TestNewModule_NameValidation(t *testing.T) {
	tests := []struct {
		name        string
		inputName   string
		expectedErr string
	}{
		{
			name:        "empty name returns error",
			inputName:   "",
			expectedErr: "name is required",
		},
		{
			name:        "whitespace only name returns error",
			inputName:   " ",
			expectedErr: "name is required",
		},
		{
			name:        "longer than 100 chars returns error",
			inputName:   "This name is really long and needs to be longer than 100 chars so we will keel typing and typiong and typ",
			expectedErr: "name cannot exceed 100 characters",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewModule(tt.inputName, "some test description")

			require.Error(t, err)
			var ve *domain.ValidationError
			require.ErrorAs(t, err, &ve)
			assert.Equal(t, tt.expectedErr, ve.Fields["name"])
		})
	}
}

func TestNewModule_DescriptionValidation(t *testing.T) {

	tests := []struct {
		name             string
		inputDescription string
		expectedErr      string
	}{
		{
			name:             "empty description returns error",
			inputDescription: "",
			expectedErr:      "description is required",
		},
		{
			name:             "whitespace only description returns error",
			inputDescription: " ",
			expectedErr:      "description is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewModule("test name", tt.inputDescription)

			require.Error(t, err)
			var ve *domain.ValidationError
			require.ErrorAs(t, err, &ve)
			assert.Equal(t, tt.expectedErr, ve.Fields["description"])

		})
	}

}

func TestNewModule_MultipleInvalidFields(t *testing.T) {
	_, err := domain.NewModule("", "")
	require.Error(t, err)
	var ve *domain.ValidationError
	require.ErrorAs(t, err, &ve)
	assert.NotEmpty(t, ve.Fields["name"])
	assert.NotEmpty(t, ve.Fields["description"])
}
