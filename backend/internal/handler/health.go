package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type HealthHandler struct {
	logger *slog.Logger
}

func NewHealthHandler(logger *slog.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

type healthResponse struct {
	Status string `json:"status"`
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	resBody, err := json.Marshal(healthResponse{Status: "ok"})
	if err != nil {
		h.logger.Error("failed to encode health response", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}
