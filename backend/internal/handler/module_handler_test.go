package handler_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
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

	body := `{"name": "Navigation Array", "description": "Controls navigation systems"}`
	r := httptest.NewRequest(http.MethodPost, "/modules", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateModule(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var response handler.CreateModuleResponse
	err := json.NewDecoder(res.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, returnedModule.ID, response.ID)
	assert.Equal(t, returnedModule.Name, response.Name)
	assert.Equal(t, returnedModule.Description, response.Description)
	assert.Equal(t, string(returnedModule.HealthState), response.HealthState)
}
