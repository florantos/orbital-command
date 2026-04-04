package handler_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/florantos/orbital-command/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestHealth_ReturnsOK(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(logger, nil, nil, nil, nil)
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	h.Health(w, r)

	res := w.Result()
	body := w.Body.String()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	assert.Equal(t, `{"status":"ok"}`, body)
}
