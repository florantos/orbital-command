package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/florantos/orbital-command/internal/domain"
	"github.com/florantos/orbital-command/internal/handler"
	"github.com/florantos/orbital-command/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCrewService struct {
	createFn func(ctx context.Context, name string, role domain.Role, qualifications []domain.Capability) (*domain.CrewMember, error)
}

func (m *mockCrewService) Create(ctx context.Context, name string, role domain.Role, qualifications []domain.Capability) (*domain.CrewMember, error) {
	return m.createFn(ctx, name, role, qualifications)
}
func TestCrewHandler_Create_Returns201OnSuccess(t *testing.T) {
	returnedCrewMember := testutil.NewTestCrewMember(t)
	crewService := &mockCrewService{
		createFn: func(ctx context.Context, name string, role domain.Role, qualifications []domain.Capability) (*domain.CrewMember, error) {
			return returnedCrewMember, nil
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(logger, nil, nil, nil, crewService)

	body := handler.CreateCrewMemberRequest{
		Name:           "John Snow",
		Role:           "engineer",
		Qualifications: []string{"docking", "hull-monitoring"},
	}

	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/crew", bytes.NewReader(bodyBytes))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateCrewMember(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var response handler.CrewResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, returnedCrewMember.ID, response.ID)
	assert.Equal(t, returnedCrewMember.Name, response.Name)
	assert.Equal(t, string(returnedCrewMember.Role), response.Role)
	expected := make([]string, len(returnedCrewMember.Qualifications))

	for i, q := range returnedCrewMember.Qualifications {
		expected[i] = string(q)
	}
	assert.Equal(t, expected, response.Qualifications)
}

func TestCrewHandler_Create_Returns400OnMalformedJSON(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(logger, nil, nil, nil, nil)

	body := `{ Name: "John Snow, Role: "engineer`

	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/crew", bytes.NewReader(bodyBytes))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateCrewMember(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var response handler.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "malformed request body", response.Error)

}

func TestCrewHandler_Create_Returns409OnDuplicateName(t *testing.T) {
	crewService := &mockCrewService{
		createFn: func(ctx context.Context, name string, role domain.Role, qualifications []domain.Capability) (*domain.CrewMember, error) {
			return nil, fmt.Errorf("create crew member: %w", domain.ErrDuplicateCrewMemberName)
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(logger, nil, nil, nil, crewService)

	body := handler.CreateCrewMemberRequest{
		Name:           "John Snow",
		Role:           "engineer",
		Qualifications: []string{"docking", "hull-monitoring"},
	}

	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/crew", bytes.NewReader(bodyBytes))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateCrewMember(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusConflict, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var response handler.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "crew member name already exists", response.Error)
}

func TestCrewHandler_Create_Returns422OnValidationFailure(t *testing.T) {
	tests := []struct {
		name           string
		mockFields     map[string]string
		expectedFields []string
	}{
		{
			name:           "single field error returns 422",
			mockFields:     map[string]string{"name": "name is required"},
			expectedFields: []string{"name"},
		},
		{
			name: "multiple field errors returns 422",
			mockFields: map[string]string{
				"name": "name is required",
				"role": "role is invalid",
			},
			expectedFields: []string{"name", "role"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crewService := &mockCrewService{
				createFn: func(ctx context.Context, name string, role domain.Role, qualifications []domain.Capability) (*domain.CrewMember, error) {
					return nil, &domain.ValidationError{Fields: tt.mockFields}
				},
			}
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			h := handler.NewHandler(logger, nil, nil, nil, crewService)

			bodyBytes, err := json.Marshal(handler.CreateCrewMemberRequest{
				Name:           "John Snow",
				Role:           "engineer",
				Qualifications: []string{"docking"},
			})
			require.NoError(t, err)

			r := httptest.NewRequest(http.MethodPost, "/crew", bytes.NewReader(bodyBytes))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.CreateCrewMember(w, r)

			res := w.Result()
			assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
			assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

			var response handler.ErrorResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, "validation failed", response.Error)

			for _, field := range tt.expectedFields {
				assert.NotEmpty(t, response.Fields[field])
			}
		})
	}
}

func TestCrewHandler_Create_Returns500OnUnexpectedError(t *testing.T) {
	crewService := &mockCrewService{
		createFn: func(ctx context.Context, name string, role domain.Role, qualifications []domain.Capability) (*domain.CrewMember, error) {
			return nil, fmt.Errorf("create crew member: unexpected database error")
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(logger, nil, nil, nil, crewService)

	body := handler.CreateCrewMemberRequest{
		Name:           "John Snow",
		Role:           "engineer",
		Qualifications: []string{"docking", "hull-monitoring"},
	}

	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/crew", bytes.NewReader(bodyBytes))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateCrewMember(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var response handler.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "internal server error", response.Error)
}
