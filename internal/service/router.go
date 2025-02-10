package service

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/recovery-flow/comtools/cifractx"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/roles"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/handlers"
	"github.com/sirupsen/logrus"
)

func Run(ctx context.Context) {
	r := chi.NewRouter()

	service, err := cifractx.GetValue[*config.Server](ctx, config.SERVER)
	if err != nil {
		logrus.Fatalf("failed to get server from context: %v", err)
	}

	r.Use(cifractx.MiddlewareWithContext(config.SERVER, service))
	authMW := service.TokenManager.AuthMdl(service.Config.JWT.AccessToken.SecretKey)
	adminGrant := service.TokenManager.RoleGrant(string(roles.RoleUserAdmin), string(roles.RoleUserSuperAdmin))

	r.Route("/re-flow", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/public", func(r chi.Router) {
				r.Post("/refresh", handlers.Refresh)

				r.Route("/oauth", func(r chi.Router) {
					r.Route("/google", func(r chi.Router) {
						r.Get("/login", handlers.GoogleLogin)
						r.Get("/callback", handlers.GoogleCallback)
					})
				})
			})

			r.Route("/private", func(r chi.Router) {
				r.Use(authMW)
				r.Route("/user", func(r chi.Router) {
					r.Route("/sessions", func(r chi.Router) {
						r.Route("/{session_id}", func(r chi.Router) {
							r.Get("/", handlers.SessionGet)
							r.Delete("/", handlers.SessionDelete)
						})

						r.Get("/", handlers.SessionsGet)
						r.Delete("/terminate", handlers.SessionsTerminate)
					})
					r.Post("/logout", handlers.Logout)
				})
			})

			r.Route("/admin", func(r chi.Router) {
				r.Use(adminGrant)
				r.Route("/{user_id}", func(r chi.Router) {
					r.Route("/sessions", func(r chi.Router) {
						r.Route("/{session_id}", func(r chi.Router) {
							r.Get("/", handlers.AdminSessionGet)
							r.Delete("/", handlers.AdminSessionDelete)
						})

						r.Get("/", handlers.AdminSessionsGet)
						r.Delete("/terminate", handlers.AdminSessionsTerminate)
					})

					r.Patch("/role/{role}", handlers.AdminRoleUpdate)
				})
			})

			r.Post("/test/login", handlers.LogSimple)
		})
	})

	server := httpkit.StartServer(ctx, service.Config.Server.Port, r, service.Logger)

	<-ctx.Done()
	httpkit.StopServer(context.Background(), server, service.Logger)
}
