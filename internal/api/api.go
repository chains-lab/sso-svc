package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hs-zavet/sso-oauth/internal/api/handlers"
	"github.com/hs-zavet/sso-oauth/internal/app"
	"github.com/hs-zavet/sso-oauth/internal/config"
	"github.com/hs-zavet/tokens"
	"github.com/hs-zavet/tokens/roles"
	"github.com/sirupsen/logrus"
)

type Api struct {
	server   *http.Server
	router   *chi.Mux
	handlers handlers.Handler

	log *logrus.Entry
	cfg config.Config
}

func NewAPI(cfg config.Config, log *logrus.Logger, app *app.App) Api {
	logger := log.WithField("module", "api")
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}
	hands := handlers.NewHandlers(cfg, logger, app)

	return Api{
		handlers: hands,
		router:   router,
		server:   server,
		log:      logger,
		cfg:      cfg,
	}
}

func (a *Api) Run(ctx context.Context, log *logrus.Logger) {
	auth := tokens.AuthMdl(a.cfg.JWT.AccessToken.SecretKey)
	admin := tokens.AccessGrant(a.cfg.JWT.AccessToken.SecretKey, roles.Admin, roles.SuperUser)

	a.router.Route("/re-news/sso", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/public", func(r chi.Router) {
				r.Post("/refresh", a.handlers.Refresh)
				r.Route("/google", func(r chi.Router) {
					r.Get("/login", a.handlers.GoogleLogin)
					r.Get("/callback", a.handlers.GoogleCallback)
				})

				r.Route("/account", func(r chi.Router) {
					r.Use(auth)
					r.Route("/sessions", func(r chi.Router) {
						r.Route("/{session_id}", func(r chi.Router) {
							r.Get("/", a.handlers.SessionGet)
							r.Delete("/", a.handlers.SessionDelete)
						})
						r.Get("/", a.handlers.SessionsGet)
						r.Delete("/", a.handlers.SessionsTerminate)
					})
					r.Get("/", a.handlers.AccountGet)
					r.Delete("/", a.handlers.Logout)
				})
			})

			r.Route("/private", func(r chi.Router) {
				r.Route("/account", func(r chi.Router) {
					r.Use(admin)
					r.Route("/{account_id}", func(r chi.Router) {
						r.Route("/sessions", func(r chi.Router) {
							r.Get("/{session_id}", a.handlers.AdminSessionGet)
							r.Get("/", a.handlers.AdminSessionsGet)
							r.Delete("/", a.handlers.AdminSessionsTerminate)
						})
						r.Get("/", a.handlers.AdminAccountGet)
						r.Patch("/{role}", a.handlers.AdminRoleUpdate)
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

func (a *Api) Start(ctx context.Context, log *logrus.Logger) {
	go func() {
		a.log.Infof("Starting server on port %s", a.cfg.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
}

func (a *Api) Stop(ctx context.Context, log *logrus.Logger) {
	a.log.Info("Shutting down server...")
	if err := a.server.Shutdown(ctx); err != nil {
		log.Errorf("Server shutdown failed: %v", err)
	}
}
