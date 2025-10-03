package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/thenopholo/shorty_url/internal/handler"
)

func Route(r *chi.Mux, h *handler.Handler) {
	// Serve static files
	fileServer := http.FileServer(http.Dir("./public"))
	r.Handle("/", fileServer)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", h.Shorten)
	})

	// Short URL redirect
	r.Get("/{code}", h.URLReciver)
}
