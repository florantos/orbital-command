package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/florantos/orbital-command/internal/domain"
)

type ModuleRepository interface {
	Create(ctx context.Context, module *domain.Module) (*domain.Module, error)
	ReadAll(ctx context.Context) ([]domain.Module, error)
}

type CreateModuleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ModuleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	HealthState string `json:"healthState"`
}

type ReadAllModulesResponse struct {
	Modules []ModuleResponse `json:"modules"`
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
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			writeValidationError(w, ve)
			return
		}
		h.logger.Error("creating module", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.logger.Info("creating module", "name", module.Name)

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

	// TODO(Loop 10): replace with authenticated user from request context
	event := domain.NewAuditEvent("module.registered", "module", created.ID, "Commander Chen", fmt.Sprintf("Registered module: %s", created.Name))

	h.logger.Info("creating audit event", "action", event.Action, "entityID", event.EntityID)

	err = h.auditEventRepo.Create(r.Context(), &h.pool, event)
	if err != nil {
		h.logger.Error("failed to create audit event", "error", err)
	} else {
		h.logger.Info("audit event created", "action", event.Action, "entityID", event.EntityID)
	}

	resp := ModuleResponse{
		ID:          created.ID,
		Name:        created.Name,
		Description: created.Description,
		HealthState: string(created.HealthState),
	}

	err = writeJSON(w, http.StatusCreated, resp)
	if err != nil {
		h.logger.Error("failed to marshal response", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
}

func (h *Handler) ReadAllModules(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("reading all modules")

	modules, err := h.moduleRepo.ReadAll(r.Context())
	if err != nil {
		h.logger.Error("failed to read all modules", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := ReadAllModulesResponse{
		Modules: []ModuleResponse{},
	}

	for _, m := range modules {
		resp.Modules = append(resp.Modules, ModuleResponse{
			ID:          m.ID,
			Name:        m.Name,
			Description: m.Description,
			HealthState: string(m.HealthState),
		})
	}
	h.logger.Info("modules read", "count", len(modules))

	err = writeJSON(w, http.StatusOK, resp)
	if err != nil {
		h.logger.Error("failed to marshal response", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

}
