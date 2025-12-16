package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func registerMiddlewares(c chi.Router) {
	c.Use(middleware.Recoverer)

	c.Use(cors.Handler(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"*"}, // AllowOrigins: "" = すべて許可
		AllowedHeaders: []string{
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
		},
		ExposedHeaders: []string{
			"Content-Disposition",
			"X-Goog-Content-Length-Range",
		},
		MaxAge: 86400, // 24時間
	}))
}
