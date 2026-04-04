package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/florantos/orbital-command/internal/domain"
)

type CrewService interface {
	Create(ctx context.Context, name string, role domain.Role, qualifications []domain.Capability) (*domain.CrewMember, error)
}

type CreateCrewMemberRequest struct {
	Name           string   `json:"name"`
	Role           string   `json:"role"`
	Qualifications []string `json:"qualifications"`
}

type CrewResponse struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Role           string   `json:"role"`
	Qualifications []string `json:"qualifications"`
}

func (h *Handler) CreateCrewMember(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed to read request body", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	var req CreateCrewMemberRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "malformed request body")
		return
	}
	quals := make([]domain.Capability, len(req.Qualifications))

	for i, q := range req.Qualifications {
		quals[i] = domain.Capability(q)
	}

	h.logger.Info("creating crew member", "name", req.Name)

	created, err := h.crewService.Create(r.Context(), req.Name, domain.Role(req.Role), quals)
	if err != nil {
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			writeValidationError(w, ve)
			return
		}
		if errors.Is(err, domain.ErrDuplicateCrewMemberName) {
			writeError(w, http.StatusConflict, "crew member name already exists")
			return
		}
		h.logger.Error("failed to create crew member", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	h.logger.Info("crew member created", "id", created.ID, "name", created.Name)

	quals2 := make([]string, len(created.Qualifications))

	for i, q := range created.Qualifications {
		quals2[i] = string(q)
	}
	resp := CrewResponse{
		ID:             created.ID,
		Name:           created.Name,
		Role:           string(created.Role),
		Qualifications: quals2,
	}

	if err := writeJSON(w, http.StatusCreated, resp); err != nil {
		h.logger.Error("failed to marshal response", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
}
