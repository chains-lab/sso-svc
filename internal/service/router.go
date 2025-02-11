package service

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/recovery-flow/comtools/cifractx"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/roles"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/handlers"
)

func Run(ctx context.Context, svc *config.Service) {
	r := chi.NewRouter()
	h := handlers.NewHandlers(svc)

	r.Use(cifractx.MiddlewareWithContext(config.SERVICE, svc))
	authMW := svc.TokenManager.AuthMdl(svc.Config.JWT.AccessToken.SecretKey)
	adminGrant := svc.TokenManager.RoleGrant(svc.Config.JWT.AccessToken.SecretKey, string(roles.RoleUserAdmin), string(roles.RoleUserSuperAdmin))

	r.Route("/re-flow", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/public", func(r chi.Router) {
				r.Post("/refresh", h.Refresh)

				r.Route("/oauth", func(r chi.Router) {
					r.Route("/google", func(r chi.Router) {
						r.Get("/login", h.GoogleLogin)
						r.Get("/callback", h.GoogleCallback)
					})
				})
			})

			r.Route("/private", func(r chi.Router) {
				r.Use(authMW)
				r.Route("/user", func(r chi.Router) {
					r.Route("/sessions", func(r chi.Router) {
						r.Route("/{session_id}", func(r chi.Router) {
							r.Get("/", h.SessionGet)
							r.Delete("/", h.SessionDelete)
						})

						r.Get("/", h.SessionsGet)
						r.Delete("/", h.SessionsTerminate)
					})
					r.Post("/logout", h.Logout)
				})
			})

			r.Route("/admin", func(r chi.Router) {
				r.Use(adminGrant)
				r.Route("/{user_id}", func(r chi.Router) {
					r.Route("/sessions", func(r chi.Router) {
						r.Route("/{session_id}", func(r chi.Router) {
							r.Get("/", h.AdminSessionGet)
							r.Delete("/", h.AdminSessionDelete)
						})

						r.Get("/", h.AdminSessionsGet)
						r.Delete("/", h.AdminSessionsTerminate)
					})

					r.Patch("/role/{role}", h.AdminRoleUpdate)
				})
			})

			r.Post("/test/login", h.LogSimple)
		})
	})

	server := httpkit.StartServer(ctx, svc.Config.Server.Port, r, svc.Logger)

	<-ctx.Done()
	httpkit.StopServer(context.Background(), server, svc.Logger)
}
