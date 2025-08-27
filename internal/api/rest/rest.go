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
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.With(userAuth).Get("/own", a.handlers.GetOwnUser)

				r.Post("/register", a.handlers.RegisterUser)
				r.Post("/login", a.handlers.Login)
				r.With(userAuth).Post("/logout", a.handlers.Logout)

				r.With(userAuth).Route("sessions", func(r chi.Router) {
					r.Get("/", a.handlers.SelectOwnSessions)
					r.Delete("/", a.handlers.DeleteOwnSessions)

					r.Route("/{session_id}", func(r chi.Router) {
						r.Get("/", a.handlers.GetOwnSession)
						r.Delete("/", a.handlers.DeleteOwnSession)
					})
				})
			})

			r.With(userAuth).Route("/admin", func(r chi.Router) {
				r.Use(adminGrant)

				r.Route("/users", func(r chi.Router) {
					r.Post("/", a.handlers.CreateUser)

					r.Route("/{user_id}", func(r chi.Router) {
						r.Get("/", a.handlers.GetUser)

						r.Route("/sessions", func(r chi.Router) {
							r.Get("/", a.handlers.SelectSessions)
							r.Delete("/", a.handlers.DeleteSessions)

							r.Route("/{session_id}", func(r chi.Router) {
								r.Get("/", a.handlers.GetSession)
								r.Delete("/", a.handlers.DeleteSession)
							})
						})
					})

					r.Put("/block", a.handlers.BlockUser)
					r.Put("/unblock", a.handlers.UnblockUser)
				})

				r.Route("/admins", func(r chi.Router) {
					r.Route("/{user_id}", func(r chi.Router) {
						r.Delete("/", a.handlers.DeleteAdmin)
					})
				})
			})
		})
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
