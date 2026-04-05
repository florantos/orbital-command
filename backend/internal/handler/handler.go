package handler

import (
	"encoding/json"
	"net/http"

	"github.com/florantos/orbital-command/internal/domain"
)

type ErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	body, err := json.Marshal(ErrorResponse{Error: message})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func writeValidationError(w http.ResponseWriter, ve *domain.ValidationError) {
	body, err := json.Marshal(ErrorResponse{Error: "validation failed", Fields: ve.Fields})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write(body)
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
	return nil
}
