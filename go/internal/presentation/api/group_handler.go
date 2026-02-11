package api

import (
	"encoding/json"
	"net/http"

	"github.com/okm321/mahking-go/internal/application"
	appin "github.com/okm321/mahking-go/internal/application/in"
	"github.com/okm321/mahking-go/pkg/logger"
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
func (h *groupHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	groups, err := h.usecase.List(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "failed to list groups", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		logger.ErrorContext(ctx, "failed to encode groups", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *groupHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input appin.CreateGroupWithRule
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	group, err := h.usecase.Create(ctx, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
		// logger.ErrorContext(ctx, "failed to create group", "error", err)
		// http.Error(w, "internal server error", http.StatusInternalServerError)
		// return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(group); err != nil {
		logger.ErrorContext(ctx, "failed to encode group", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
