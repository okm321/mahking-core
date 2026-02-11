package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/okm321/mahking-go/pkg/logger"
	"github.com/riandyrn/otelchi"
)

// httpLogger はCloud Loggingでは不要なため一旦未使用
//
//nolint:unused
func httpLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()

		defer func() {
			status := ww.Status()
			if status == 0 {
				status = http.StatusOK
			}

			logger.InfoContext(r.Context(), "http request",
				logger.HTTPAttr(r, status, time.Since(start), ww.BytesWritten()),
			)
		}()

		next.ServeHTTP(ww, r)
	})
}

func registerMiddlewares(c chi.Router) {
	c.Use(middleware.Recoverer)
	c.Use(otelchi.Middleware("mahking-go", otelchi.WithChiRoutes(c)))

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
