package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Instance *http.Server
}

func NewServer(port string, router *chi.Mux) *Server {
	return &Server{
		Instance: &http.Server{
			Handler: router,
			Addr:    ":" + port,
		},
	}
}

func (s *Server) Start() {
	slog.Info("starting server on port " + s.Instance.Addr)
	if err := s.Instance.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("error starting server", "err", err)
	}
}

func (s *Server) Closer(ctx context.Context) error {
	slog.Info("shutting down server...")
	return s.Instance.Shutdown(ctx)
}
