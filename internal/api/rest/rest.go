package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal/api/rest/handlers"
	"github.com/chains-lab/sso-svc/internal/api/rest/mdlv"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/go-chi/chi/v5"
)

type Rest struct {
	server   *http.Server
	router   *chi.Mux
	handlers handlers.Service

	log logium.Logger
	cfg config.Config
}

func NewRest(cfg config.Config, log logium.Logger, app *app.App) Rest {
	logger := log.WithField("module", "api")
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}
	hands := handlers.NewService(cfg, logger, app)

	return Rest{
		handlers: hands,
		router:   router,
		server:   server,
		log:      logger,
		cfg:      cfg,
	}
}

func (a *Rest) Run(ctx context.Context) {
	userAuth := mdlv.AuthMdl(a.cfg.JWT.User.AccessToken.SecretKey)
	adminGrant := mdlv.AccessGrant(roles.Admin, roles.SuperUser)
	svcAuth := mdlv.ServiceAuthMdl(a.cfg.JWT.Service.SecretKey)

	a.log.WithField("module", "api").Info("Starting API server")

	a.router.Route("/sso-svc/", func(r chi.Router) {
		r.Use(svcAuth)
	})

	a.Start(ctx)

	<-ctx.Done()
	a.Stop(ctx)
}

func (a *Rest) Start(ctx context.Context) {
	go func() {
		a.log.Infof("Starting server on port %s", a.cfg.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Fatalf("Server failed to start: %v", err)
		}
	}()
}

func (a *Rest) Stop(ctx context.Context) {
	a.log.Info("Shutting down server...")
	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Errorf("Server shutdown failed: %v", err)
	}
}
