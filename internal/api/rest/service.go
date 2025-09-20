package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/go-chi/chi/v5"
)

type Service struct {
	server *http.Server
	router *chi.Mux
	log    logium.Logger
	cfg    config.Config
}

func NewRest(cfg config.Config, log logium.Logger) Service {
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}

	router.Use()

	return Service{
		router: router,
		server: server,
		log:    log,
		cfg:    cfg,
	}
}

func (s *Service) Start(ctx context.Context) {
	go func() {
		s.log.Infof("Starting server on port %s", s.cfg.Server.Port)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Fatalf("Server failed to start: %v", err)
		}
	}()
}

func (s *Service) Stop(ctx context.Context) {
	s.log.Info("Shutting down server...")
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Errorf("Server shutdown failed: %v", err)
	}
}
