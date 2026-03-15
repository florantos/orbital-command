package handler

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	resBody, err := json.Marshal(healthResponse{Status: "ok"})
	if err != nil {
		h.logger.Error("failed to encode health response", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}
