package handler_test

import (
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

type mockAuditEventService struct {
	readAllFn func(ctx context.Context) ([]domain.AuditEvent, error)
}

func (m *mockAuditEventService) ReadAll(ctx context.Context) ([]domain.AuditEvent, error) {
	return m.readAllFn(ctx)
}

func TestAuditEventHandler_ReadAll_Returns200(t *testing.T) {
	tests := []struct {
		name  string
		input []domain.AuditEvent
	}{
		{
			name: "returns audit events on success",
			input: []domain.AuditEvent{
				*testutil.NewTestAuditEvent(t),
				*testutil.NewTestAuditEvent(t),
			},
		},
		{
			name:  "returns empty slice on success",
			input: []domain.AuditEvent{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			auditEventService := &mockAuditEventService{
				readAllFn: func(ctx context.Context) ([]domain.AuditEvent, error) {
					return tt.input, nil
				},
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			h := handler.NewAuditHandler(logger, auditEventService)

			r := httptest.NewRequest(http.MethodGet, "/audit-events", nil)
			w := httptest.NewRecorder()

			h.ReadAllAuditEvents(w, r)

			res := w.Result()
			assert.Equal(t, http.StatusOK, res.StatusCode)

			var response handler.ReadAllAuditEventsResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.NotNil(t, response.AuditEvents)
			assert.Len(t, response.AuditEvents, len(tt.input))
		})
	}

}
func TestAuditEventHandler_ReadAll_Returns500OnUnexpectedError(t *testing.T) {
	auditEventService := &mockAuditEventService{
		readAllFn: func(ctx context.Context) ([]domain.AuditEvent, error) {
			return []domain.AuditEvent{}, fmt.Errorf("read all audit events: unexpected database error")
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewAuditHandler(logger, auditEventService)

	r := httptest.NewRequest(http.MethodGet, "/audit-events", nil)
	w := httptest.NewRecorder()

	h.ReadAllAuditEvents(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var response handler.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "internal server error", response.Error)
}
