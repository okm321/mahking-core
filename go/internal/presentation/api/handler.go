package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

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
		logger.WarnContext(ctx, err.Error(), "error", err)
		writeJSON(w, http.StatusBadRequest, pkgerror.ErrorResponse{
			Message: validationMessage(validationErrs),
			Reason:  string(pkgerror.ErrCodeValidation),
		})
		return
	}

	var pkgErr *pkgerror.Error
	if errors.As(err, &pkgErr) {
		logger.WarnContext(ctx, err.Error(), "error", err)
		writeJSON(w, http.StatusBadRequest, pkgErr)
		return
	}

	logger.ErrorContext(ctx, err.Error(), "error", err)
	writeJSON(w, http.StatusInternalServerError, pkgerror.ErrorResponse{
		Message: "internal server error",
		Reason:  string(pkgerror.ErrCodeInternal),
	})
}

func validationMessage(errs govaliderrors.ValidationErrors) string {
	msgs := make([]string, 0, len(errs))
	for _, e := range errs {
		msgs = append(msgs, e.Reason)
	}
	return strings.Join(msgs, ", ")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
