package server

import (
	"github.com/go-chi/chi"
	"github.com/thenopholo/shorty_url/internal/handler"
)

func Route(r *chi.Mux, h *handler.Handler) {
	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", h.Shorten)
		r.Get("/{code}", h.URLReciver)
	})
}
