package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/thenopholo/shorty_url/internal/service"
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

func sendJSON(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("failed to marshal json data", "error", err)
		sendJSON(w, Response{Err: "something went wrong"}, http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write response to client", "error", err)
		return
	}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var body RequestBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sendJSON(w, Response{Err: "invalid body"}, http.StatusUnprocessableEntity)
		return
	}

	if _, err := url.Parse(body.URL); err != nil {
		sendJSON(w, Response{Err: "invalid url"}, http.StatusBadRequest)
		return
	}

	var code string
	var err error
	maxAttempts := 5
	for i := 0; i < maxAttempts; i++ {
		code, err = service.GenerateRandomCode(6)
		if err != nil {
			sendJSON(w, Response{Err: "failed to generate code"}, http.StatusInternalServerError)
			return
		}

		_, err = h.store.GetURL(code)
		if err != nil {
			break
		}
	}

	if err := h.store.SaveURL(code, body.URL); err != nil {
		sendJSON(w, Response{Err: "failed to save url"}, http.StatusInternalServerError)
		return
	}

	sendJSON(w, Response{Data: map[string]string{"code": code}}, http.StatusCreated)
}

func (h *Handler) URLReciver(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	url, err := h.store.GetURL(code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendJSON(w, Response{Err: "original url not found"}, http.StatusNotFound)
			return
		}
		sendJSON(w, Response{Err: "something went wrong"}, http.StatusInternalServerError)
		return }

	http.Redirect(w, r, url, http.StatusPermanentRedirect)

}
