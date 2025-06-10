package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/app"
	"github.com/chains-lab/chains-auth/internal/config"
	"github.com/chains-lab/chains-auth/internal/rest/handlers"
	"github.com/chains-lab/gatekit/mdlv"
	"github.com/chains-lab/gatekit/roles"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Rest struct {
	server   *http.Server
	router   *chi.Mux
	handlers handlers.Handlers

	log *logrus.Entry
	cfg config.Config
}

func NewRest(cfg config.Config, log *logrus.Logger, app *app.App) Rest {
	logger := log.WithField("module", "api")
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}
	hands := handlers.NewHandlers(cfg, logger, app)

	return Rest{
		handlers: hands,
		router:   router,
		server:   server,
		log:      logger,
		cfg:      cfg,
	}
}

func (a *Rest) Run(ctx context.Context, log *logrus.Logger) {
	auth := mdlv.AuthMdl(a.cfg.JWT.AccessToken.SecretKey, "todo")
	admin := mdlv.AccessGrant(a.cfg.JWT.AccessToken.SecretKey, "todo", roles.Admin, roles.SuperUser)

	a.log.WithField("module", "api").Info("Starting API server")

	a.router.Route("/chains/auth", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

			r.Route("/own", func(r chi.Router) {
				r.Use(auth)

				r.Post("/refresh", a.handlers.Refresh)

				r.Get("/login", a.handlers.GoogleLogin)
				r.Get("/callback", a.handlers.GoogleCallback)

				r.Delete("/logout", a.handlers.Logout)

				r.Get("/", a.handlers.OwnUserGet)

				r.Route("/sessions", func(r chi.Router) {
					r.Route("/{session_id}", func(r chi.Router) {
						r.Get("/", a.handlers.OwnGetSession)
						r.Delete("/", a.handlers.DeleteSession)
					})

					r.Get("/", a.handlers.OwnGetSessions)
					r.Delete("/", a.handlers.OwnTerminateSessions)
				})
			})

			r.Route("/admin", func(r chi.Router) {
				r.Use(admin)

				r.Route("/{user_id}", func(r chi.Router) {
					r.Get("/", a.handlers.AdminGetUser)
					r.Patch("/{role}", a.handlers.AdminUpdateRole)

					r.Route("/sessions", func(r chi.Router) {
						r.Route("/{session_id}", func(r chi.Router) {
							r.Get("/", a.handlers.AdminGetSession)
							r.Delete("/", a.handlers.AdminDeleteSession)
						})
						r.Get("/", a.handlers.AdminGetSessions)
						r.Delete("/", a.handlers.AdminTerminateSessions)
					})
				})
			})

			r.Post("/test/login", a.handlers.LoginSimple)
		})
	})

	a.Start(ctx, log)

	<-ctx.Done()
	a.Stop(ctx, log)
}

func (a *Rest) Start(ctx context.Context, log *logrus.Logger) {
	go func() {
		a.log.Infof("Starting server on port %s", a.cfg.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
}

func (a *Rest) Stop(ctx context.Context, log *logrus.Logger) {
	a.log.Info("Shutting down server...")
	if err := a.server.Shutdown(ctx); err != nil {
		log.Errorf("Server shutdown failed: %v", err)
	}
}
