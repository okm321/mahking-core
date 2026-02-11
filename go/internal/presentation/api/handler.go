package api

import (
	"encoding/json"
	"errors"
	"net/http"

	pkgerror "github.com/okm321/mahking-go/pkg/error"
	"github.com/okm321/mahking-go/pkg/logger"
	govaliderrors "github.com/sivchari/govalid/validation/errors"
)

// AppHandler は error を返せるハンドラー型
type AppHandler func(w http.ResponseWriter, r *http.Request) error

// Handle は AppHandler を http.HandlerFunc に変換するアダプタ
func Handle(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			handleError(w, r, err)
		}
	}
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()

	var notFoundErr *pkgerror.ErrorNotFound
	if errors.As(err, &notFoundErr) {
		logger.WarnContext(ctx, "not found", "error", err)
		writeJSON(w, http.StatusNotFound, notFoundErr)
		return
	}

	var validationErrs govaliderrors.ValidationErrors
	if errors.As(err, &validationErrs) {
		logger.ErrorContext(ctx, "validation error", "error", err)
		writeJSON(w, http.StatusBadRequest, validationErrs)
		return
	}

	var pkgErr *pkgerror.Error
	if errors.As(err, &pkgErr) {
		logger.WarnContext(ctx, "bad request", "error", err)
		writeJSON(w, http.StatusBadRequest, pkgErr)
		return
	}

	logger.ErrorContext(ctx, "internal server error", "error", err)
	writeJSON(w, http.StatusInternalServerError, map[string]string{
		"message": "internal server error",
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
