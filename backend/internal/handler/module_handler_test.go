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

	"github.com/florantos/orbital-command/internal/database"
	"github.com/florantos/orbital-command/internal/domain"
	"github.com/florantos/orbital-command/internal/handler"
	"github.com/florantos/orbital-command/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockModuleRepo struct {
	createFn  func(ctx context.Context, module *domain.Module) (*domain.Module, error)
	readAllFn func(ctx context.Context) ([]domain.Module, error)
}

func (m *mockModuleRepo) Create(ctx context.Context, module *domain.Module) (*domain.Module, error) {
	return m.createFn(ctx, module)
}

func (m *mockModuleRepo) ReadAll(ctx context.Context) ([]domain.Module, error) {
	return m.readAllFn(ctx)
}

func TestModuleHandler_Create_Returns201OnSuccess(t *testing.T) {
	returnedModule := testutil.NewTestModule(t)

	moduleRepo := &mockModuleRepo{
		createFn: func(ctx context.Context, module *domain.Module) (*domain.Module, error) {
			return returnedModule, nil
		},
	}
	auditEventRepo := &mockAuditEventRepo{
		createFn: func(ctx context.Context, db database.DBTX, event *domain.AuditEvent) error {
			return nil
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(logger, nil, moduleRepo, auditEventRepo, nil)

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

	var response handler.ModuleResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, returnedModule.ID, response.ID)
	assert.Equal(t, returnedModule.Name, response.Name)
	assert.Equal(t, returnedModule.Description, response.Description)
	assert.Equal(t, string(returnedModule.HealthState), response.HealthState)
}

func TestModuleHandler_Create_Returns400OnMalformedJSON(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(logger, nil, nil, nil, nil)

	body := `{ Name: "Navigation Array, Description: "Controls navigation systems`

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

func TestModuleHandler_Create_Returns409OnDuplicateName(t *testing.T) {
	moduleRepo := &mockModuleRepo{
		createFn: func(ctx context.Context, module *domain.Module) (*domain.Module, error) {
			return nil, fmt.Errorf("create module: %w", domain.ErrDuplicateModuleName)
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(logger, nil, moduleRepo, nil, nil)

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

func TestModuleHandler_Create_Returns500OnUnexpectedError(t *testing.T) {
	moduleRepo := &mockModuleRepo{
		createFn: func(ctx context.Context, module *domain.Module) (*domain.Module, error) {
			return nil, fmt.Errorf("create module: unexpected database error")
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(logger, nil, moduleRepo, nil, nil)

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

func TestModuleHandler_Create_Returns422OnValidationFailure(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        handler.CreateModuleRequest
		expectedFields []string
	}{
		{
			name: "name empty returns 422",
			reqBody: handler.CreateModuleRequest{
				Name:        "",
				Description: "Controls navigation systems",
			},
			expectedFields: []string{"name"},
		},
		{
			name: "whitespace only name returns 422",
			reqBody: handler.CreateModuleRequest{
				Name:        " ",
				Description: "Controls navigation systems",
			},
			expectedFields: []string{"name"},
		},
		{
			name: "name longer than 100 chars returns 422",
			reqBody: handler.CreateModuleRequest{
				Name:        "This name is really long and needs to be longer than 100 chars so we will keel typing and typiong and typ",
				Description: "Controls navigation systems",
			},
			expectedFields: []string{"name"},
		},
		{
			name: "description empty returns 422",
			reqBody: handler.CreateModuleRequest{
				Name:        "Navigation Array",
				Description: "",
			},
			expectedFields: []string{"description"},
		},
		{
			name: "whitespace only description returns 422",
			reqBody: handler.CreateModuleRequest{
				Name:        "Navigation Array",
				Description: " ",
			},
			expectedFields: []string{"description"},
		},
		{
			name: "multiple field errors returns 422",
			reqBody: handler.CreateModuleRequest{
				Name:        "",
				Description: "",
			},
			expectedFields: []string{"name", "description"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			h := handler.NewHandler(logger, nil, nil, nil, nil)

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
			assert.Equal(t, "validation failed", response.Error)

			for _, field := range tt.expectedFields {
				assert.NotEmpty(t, response.Fields[field])

			}
		})
	}
}

func TestModuleHandler_ReadAll_Returns500OnUnexpectedError(t *testing.T) {
	moduleRepo := &mockModuleRepo{
		readAllFn: func(ctx context.Context) ([]domain.Module, error) {
			return []domain.Module{}, fmt.Errorf("read all modules: unexpected database error")
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(logger, nil, moduleRepo, nil, nil)

	r := httptest.NewRequest(http.MethodGet, "/modules", nil)
	w := httptest.NewRecorder()

	h.ReadAllModules(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var response handler.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "internal server error", response.Error)

}

func TestModulesHandler_ReadAll_Returns200(t *testing.T) {
	tests := []struct {
		name  string
		input []domain.Module
	}{
		{
			name: "returns modules on success",
			input: []domain.Module{
				*testutil.NewTestModule(t),
				*testutil.NewTestModule(t),
			},
		},
		{
			name:  "returns empty slice on success",
			input: []domain.Module{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moduleRepo := &mockModuleRepo{
				readAllFn: func(ctx context.Context) ([]domain.Module, error) {
					return tt.input, nil
				},
			}
			auditEventRepo := &mockAuditEventRepo{
				createFn: func(ctx context.Context, db database.DBTX, event *domain.AuditEvent) error {
					return nil
				},
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			h := handler.NewHandler(logger, nil, moduleRepo, auditEventRepo, nil)

			r := httptest.NewRequest(http.MethodGet, "/modules", nil)
			w := httptest.NewRecorder()

			h.ReadAllModules(w, r)

			res := w.Result()
			assert.Equal(t, http.StatusOK, res.StatusCode)

			var response handler.ReadAllModulesResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.NotNil(t, response.Modules)
			assert.Len(t, response.Modules, len(tt.input))
		})
	}

}

func TestModuleHandler_Create_EmitsAuditEventOnSuccess(t *testing.T) {
	returnedModule := testutil.NewTestModule(t)

	moduleRepo := &mockModuleRepo{
		createFn: func(ctx context.Context, module *domain.Module) (*domain.Module, error) {
			return returnedModule, nil
		},
	}

	auditEventRepo := &mockAuditEventRepo{
		createFn: func(ctx context.Context, db database.DBTX, event *domain.AuditEvent) error {
			return nil
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(logger, nil, moduleRepo, auditEventRepo, nil)

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
	assert.True(t, auditEventRepo.called)

	var response handler.ModuleResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, returnedModule.ID, response.ID)
	assert.Equal(t, returnedModule.Name, response.Name)
	assert.Equal(t, returnedModule.Description, response.Description)
	assert.Equal(t, string(returnedModule.HealthState), response.HealthState)
}
