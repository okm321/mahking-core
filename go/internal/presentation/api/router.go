package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/okm321/mahking-go/pkg/logger"
)

// Router is an alias to make dependency injection explicit.
type Router = chi.Router

type HandlerSet struct {
	Group *groupHandler
}

// NewRouter sets up middleware and routes.
func NewRouter(handlers HandlerSet) Router {
	r := chi.NewRouter()

	registerMiddlewares(r)

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		logger.WarnContext(context.Background(), "root endpoint accessed")
		_, _ = w.Write([]byte("welcome"))
	})

	r.Route("/groups", func(r chi.Router) {
		r.Get("/", Handle(handlers.Group.List))
		r.Post("/", Handle(handlers.Group.Create))
	})

	return r
}
