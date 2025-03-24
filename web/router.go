package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	SetupMiddleware(r)

	// Serve static files
	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/*", fileServer)

	r.Get("/translate", TranslatePageHandler)

	// API endpoints
	r.Route("/api", func(r chi.Router) {
		r.Get("/languages", GetLanguagesHandler)
		r.Get("/translate", TranslateHandler)
	})

	return r
}
