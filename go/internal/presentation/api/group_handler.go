package api

import (
	"encoding/json"
	"net/http"

	"github.com/okm321/mahking-go/internal/application"
	appin "github.com/okm321/mahking-go/internal/application/in"
	pkgerror "github.com/okm321/mahking-go/pkg/error"
)

type groupHandler struct {
	usecase *application.GroupUsecase
}

func NewGroupHandler(usecase *application.GroupUsecase) *groupHandler {
	return &groupHandler{
		usecase: usecase,
	}
}

// List returns all groups.
func (h *groupHandler) List(w http.ResponseWriter, r *http.Request) error {
	groups, err := h.usecase.List(r.Context())
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(groups)
}

func (h *groupHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var input appin.CreateGroupWithRule
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return pkgerror.NewError("invalid json body")
	}
	group, err := h.usecase.Create(r.Context(), input)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(group)
}
