package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thenopholo/shorty_url/internal/handler"
	"github.com/thenopholo/shorty_url/internal/store"
)

type Server struct {
	port    string
	mux     *chi.Mux
	server  *http.Server
	logger  *slog.Logger
	handler *handler.Handler
}

type Config struct {
	Port   string
	Logger *slog.Logger
	Store  *store.Store
}

func NewServer(cfg Config) *Server {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.RequestID)

  h := handler.NewHandler(cfg.Store)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		port:   srv.Addr,
		mux:    mux,
		server: srv,
		logger: cfg.Logger,
    handler: h,
	}
}

func (s *Server) SetupRoutes() {
	Route(s.mux, s.handler)
}

func (s *Server) Start() error {
	s.logger.Info("Starting server", "port", s.port)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server")
	return s.server.Shutdown(ctx)
}
