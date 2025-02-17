package rest

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/roles"
	"github.com/recovery-flow/sso-oauth/internal/service"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/rest/handlers"
)

func Run(ctx context.Context, svc *service.Service) {
	r := chi.NewRouter()
	h := handlers.NewHandlers(svc)

	authMW := svc.TokenManager.AuthMdl(svc.Config.JWT.AccessToken.SecretKey)
	adminGrant := svc.TokenManager.RoleGrant(svc.Config.JWT.AccessToken.SecretKey, string(roles.RoleUserAdmin), string(roles.RoleUserSuperAdmin))

	r.Route("/sso", func(r chi.Router) {
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
				r.Route("/account", func(r chi.Router) {
					r.Get("/", h.AccountGet)

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
				r.Route("/{account_id}", func(r chi.Router) {
					r.Route("/sessions", func(r chi.Router) {
						r.Get("/", h.AdminSessionsGet)
						r.Delete("/", h.AdminSessionsTerminate)
					})

					r.Patch("/role/{role}", h.AdminRoleUpdate)
				})

				r.Route("/sessions", func(r chi.Router) {
					r.Route("/{session_id}", func(r chi.Router) {
						r.Get("/", h.AdminSessionsGet)
						r.Delete("/", h.AdminSessionsTerminate)
					})
				})
			})

			r.Post("/test/login", h.LoginSimple)
		})
	})

	server := httpkit.StartServer(ctx, svc.Config.Server.Port, r, svc.Logger)

	<-ctx.Done()
	httpkit.StopServer(context.Background(), server, svc.Logger)
}
