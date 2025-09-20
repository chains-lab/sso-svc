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

func (a *Service) Start(ctx context.Context) {
	go func() {
		a.log.Infof("Starting server on port %s", a.cfg.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Fatalf("Server failed to start: %v", err)
		}
	}()
}

func (a *Service) Stop(ctx context.Context) {
	a.log.Info("Shutting down server...")
	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Errorf("Server shutdown failed: %v", err)
	}
}
