package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Health check
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Placeholder routes for future endpoints
		r.Route("/products", func(r chi.Router) {
			r.Get("/", http.NotFound)
			r.Get("/{id}", http.NotFound)
		})
		r.Route("/competitors", func(r chi.Router) {
			r.Get("/", http.NotFound)
			r.Get("/{id}", http.NotFound)
		})
		r.Route("/trends", func(r chi.Router) {
			r.Get("/", http.NotFound)
		})
		r.Route("/strategy", func(r chi.Router) {
			r.Get("/pricing", http.NotFound)
		})
	})

	return r
}
