package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/florantos/orbital-command/internal/domain"
)

type ModuleRepository interface {
	Create(ctx context.Context, module *domain.Module) (*domain.Module, error)
}

type CreateModuleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateModuleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	HealthState string `json:"healthState"`
}

func (h *Handler) CreateModule(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed to read request body", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	var req CreateModuleRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "malformed request body")
		return
	}

	module, err := domain.NewModule(req.Name, req.Description)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	h.logger.Info("creating module", "name", req.Name)

	created, err := h.moduleRepo.Create(r.Context(), module)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateModuleName) {
			writeError(w, http.StatusConflict, "module name already exists")
			return
		}
		h.logger.Error("failed to create module", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.logger.Info("module created", "id", created.ID, "name", created.Name)

	resp := CreateModuleResponse{
		ID:          created.ID,
		Name:        created.Name,
		Description: created.Description,
		HealthState: string(created.HealthState),
	}

	resBody, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error("failed to marshal response", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resBody) //nolint:errcheck
}
