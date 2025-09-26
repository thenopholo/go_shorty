package handler

import (
	"net/http"

	"github.com/thenopholo/shorty_url/internal/store"
)

type Handler struct {
	store *store.Store
}

func NewHandler(s *store.Store) *Handler {
	return &Handler{
		store: s,
	}
}

type RequestBody struct {
	URL string `json:"url"`
}

type Response struct {
	Err  string `json:"error,omitempty"`
	Data any    `json:"data,omitempty"`
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) URLReciver(w http.ResponseWriter, r *http.Request) {

}
