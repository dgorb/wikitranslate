package web

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupMiddleware(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// TODO: Look into re-enabling.
	// r.Use(CSPMiddleware)
}

func CSPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' 'unsafe-inline' cdn.tailwindcss.com; "+
				"style-src 'self' 'unsafe-inline' kit.fontawesome.com; "+
				"style-src 'self' 'unsafe-inline' fonts.googleapis.com; "+
				"font-src fonts.gstatic.com;")

		next.ServeHTTP(w, r)
	})
}
