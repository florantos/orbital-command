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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockModuleRepo struct {
	createFn func(ctx context.Context, module *domain.Module) (*domain.Module, error)
}

func (m *mockModuleRepo) Create(ctx context.Context, module *domain.Module) (*domain.Module, error) {
	return m.createFn(ctx, module)
}

func TestCreateModuleHandler_Returns201OnSuccess(t *testing.T) {
	returnedModule := &domain.Module{
		ID:          "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		Name:        "Navigation Array",
		Description: "Controls navigation systems",
		HealthState: domain.HealthStateOperational,
	}

	moduleRepo := &mockModuleRepo{
		createFn: func(ctx context.Context, module *domain.Module) (*domain.Module, error) {
			return returnedModule, nil
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(logger, moduleRepo)

	body := handler.CreateModuleRequest{
		Name:        "Navigation Array",
		Description: "Controls navigation systems",
	}
	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/modules", bytes.NewReader(bodyBytes))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateModule(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var response handler.CreateModuleResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, returnedModule.ID, response.ID)
	assert.Equal(t, returnedModule.Name, response.Name)
	assert.Equal(t, returnedModule.Description, response.Description)
	assert.Equal(t, string(returnedModule.HealthState), response.HealthState)
}

func TestCreateModuleHandler_Returns400OnMalformedJSON(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(logger, nil)

	body := handler.CreateModuleRequest{
		Name:        "Navigation Array",
		Description: "Controls navigation systems",
	}
	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/modules", bytes.NewReader(bodyBytes))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateModule(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var response handler.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "malformed request body", response.Error)
}

func TestCreateModuleHandler_Returns409OnDuplicateName(t *testing.T) {
	moduleRepo := &mockModuleRepo{
		createFn: func(ctx context.Context, module *domain.Module) (*domain.Module, error) {
			return nil, fmt.Errorf("create module: %w", domain.ErrDuplicateModuleName)
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(logger, moduleRepo)

	body := handler.CreateModuleRequest{
		Name:        "Navigation Array",
		Description: "Controls navigation systems",
	}
	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/modules", bytes.NewReader(bodyBytes))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateModule(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusConflict, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var response handler.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "module name already exists", response.Error)

}

func TestCreateModuleHandler_Returns500OnUnexpectedError(t *testing.T) {
	moduleRepo := &mockModuleRepo{
		createFn: func(ctx context.Context, module *domain.Module) (*domain.Module, error) {
			return nil, fmt.Errorf("create module: unexpected database error")
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(logger, moduleRepo)

	body := handler.CreateModuleRequest{
		Name:        "Navigation Array",
		Description: "Controls navigation systems",
	}
	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/modules", bytes.NewReader(bodyBytes))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateModule(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var response handler.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "internal server error", response.Error)
}

func TestCreateModuleHandler_Returns422OnValidationFailure(t *testing.T) {
	tests := []struct {
		name        string
		reqBody     handler.CreateModuleRequest
		expectedErr string
	}{
		{
			name: "name empty returns 422",
			reqBody: handler.CreateModuleRequest{
				Name:        "",
				Description: "Controls navigation systems",
			},
			expectedErr: "name is required",
		},
		{
			name: "whitespace only name returns 422",
			reqBody: handler.CreateModuleRequest{
				Name:        " ",
				Description: "Controls navigation systems",
			},
			expectedErr: "name is required",
		},
		{
			name: "name longer than 100 chars returns 422",
			reqBody: handler.CreateModuleRequest{
				Name:        "This name is really long and needs to be longer than 100 chars so we will keel typing and typiong and typ",
				Description: "Controls navigation systems",
			},
			expectedErr: "name cannot exceed 100 characters",
		},
		{
			name: "description empty returns 422",
			reqBody: handler.CreateModuleRequest{
				Name:        "Navigation Array",
				Description: "",
			},
			expectedErr: "description is required",
		},
		{
			name: "whitepase only description returns 422",
			reqBody: handler.CreateModuleRequest{
				Name:        "Navigation Array",
				Description: " ",
			},
			expectedErr: "description is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			h := handler.NewHandler(logger, nil)

			bodyBytes, err := json.Marshal(tt.reqBody)
			require.NoError(t, err)

			r := httptest.NewRequest(http.MethodPost, "/modules", bytes.NewReader(bodyBytes))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.CreateModule(w, r)

			res := w.Result()
			assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
			assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

			var response handler.ErrorResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedErr, response.Error)
		})
	}
}

func TestCreateMOduleHandler_EmitsAuditEventOnSuccess(t *testing.T) {

}
