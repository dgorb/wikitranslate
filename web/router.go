package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)

	// Serve static files
	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/", fileServer)
	r.Handle("/scripts.js", fileServer)

	// API endpoints
	r.Get("/languages", GetLanguagesHandler)
	r.Get("/translate", TranslateHandler)

	return r
}
