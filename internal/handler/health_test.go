package handler_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/florantos/orbital-command/internal/handler"
)

func TestHealth_ReturnsOK(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(logger)

	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	h.Health(w, r)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	if res.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", res.Header.Get("Content-Type"))
	}

	body := w.Body.String()
	expected := `{"status":"ok"}`
	if body != expected {
		t.Errorf("expected body %q, got %q", expected, body)
	}
}
