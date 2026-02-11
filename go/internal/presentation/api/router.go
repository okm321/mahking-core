package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Router is an alias to make dependency injection explicit.
type Router = chi.Router

type HandlerSet struct {
	Group *groupHandler
}

// NewRouter sets up middleware and routes.
func NewRouter(handlers HandlerSet) Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	registerMiddlewares(r)

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("welcome"))
	})

	r.Route("/groups", func(r chi.Router) {
		r.Get("/", handlers.Group.List)
		r.Post("/", handlers.Group.Create)
	})

	return r
}
