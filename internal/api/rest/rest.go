package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/chains-lab/gatekit/mdlv"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal/api/rest/handlers"
	"github.com/chains-lab/sso-svc/internal/api/rest/meta"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/go-chi/chi/v5"
)

type Rest struct {
	server *http.Server
	router *chi.Mux

	handlers handlers.Service
	log      logium.Logger
	cfg      config.Config
}

func NewRest(cfg config.Config, log logium.Logger, app *app.App) Rest {
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}
	hands := handlers.NewService(cfg, log, app)

	router.Use()

	return Rest{
		handlers: hands,
		router:   router,
		server:   server,
		log:      log,
		cfg:      cfg,
	}
}

func (a *Rest) Run(ctx context.Context) {
	//svcAuth := mdlv.ServiceAuthMdl(enum.ServiceName, a.cfg.JWT.Service.SecretKey)
	userAuth := mdlv.AuthMdl(meta.UserCtxKey, a.cfg.JWT.User.AccessToken.SecretKey)
	adminGrant := mdlv.AccessGrant(meta.UserCtxKey, roles.Admin, roles.SuperUser)

	a.router.Route("/sso-svc/", func(r chi.Router) {
		//r.Use(svcAuth)
		r.Route("/v1", func(r chi.Router) {
			r.Post("/register", a.handlers.RegisterUser)

			r.Route("/login", func(r chi.Router) {
				r.Post("/", a.handlers.Login)

				r.Route("/google", func(r chi.Router) {
					r.Post("/", a.handlers.GoogleLogin)
					r.Post("/callback", a.handlers.GoogleCallback)
				})
			})

			r.Post("/refresh", a.handlers.RefreshToken)

			r.With(userAuth).Route("/own", func(r chi.Router) {
				r.With(userAuth).Get("/", a.handlers.GetOwnUser)
				r.With(userAuth).Post("/logout", a.handlers.Logout)
				r.With(userAuth).Post("/password", a.handlers.UpdatePassword)

				r.With(userAuth).Route("/sessions", func(r chi.Router) {
					r.Get("/", a.handlers.SelectOwnSessions)
					r.Delete("/", a.handlers.DeleteOwnSessions)

					r.Route("/{session_id}", func(r chi.Router) {
						r.Get("/", a.handlers.GetOwnSession)
						r.Delete("/", a.handlers.DeleteOwnSession)
					})
				})
			})

			r.Route("/admin", func(r chi.Router) {
				r.Use(userAuth)
				r.Use(adminGrant)

				r.Post("/", a.handlers.RegisterAdmin)

				r.Route("/{user_id}", func(r chi.Router) {
					r.Get("/", a.handlers.GetUser)

					r.Route("/sessions", func(r chi.Router) {
						r.Get("/", a.handlers.SelectUserSessions)
						r.Delete("/", a.handlers.DeleteSessions)

						r.Route("/{session_id}", func(r chi.Router) {
							r.Get("/", a.handlers.GetSession)
							r.Delete("/", a.handlers.DeleteSession)
						})
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
