package api

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/sso-oauth/internal/service"
	"github.com/recovery-flow/sso-oauth/internal/service/api/handlers"
	"github.com/recovery-flow/tokens"
	"github.com/recovery-flow/tokens/identity"
)

func Run(ctx context.Context, svc *service.Service) {
	r := chi.NewRouter()

	h, err := handlers.NewHandlers(svc)
	if err != nil {
		svc.Log.Fatalf("failed to create handlers: %v", err)
		<-ctx.Done()
		return
	}

	authMW := tokens.AuthMdl(svc.Config.JWT.AccessToken.SecretKey)
	roleGrant := tokens.IdentityMdl(svc.Config.JWT.AccessToken.SecretKey, identity.Admin, identity.SuperUser)

	r.Route("/re-flow/sso", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/public", func(r chi.Router) {
				r.Post("/refresh", h.Refresh)
				r.Route("/google", func(r chi.Router) {
					r.Get("/login", h.GoogleLogin)
					r.Get("/callback", h.GoogleCallback)
				})
			})

			r.Route("/private", func(r chi.Router) {
				r.Route("/account", func(r chi.Router) {
					r.With(authMW).Route("/current", func(r chi.Router) {
						r.Route("/sessions", func(r chi.Router) {
							r.Route("/{session_id}", func(r chi.Router) {
								r.Get("/", h.SessionGet)
								r.Delete("/", h.SessionDelete)
							})
							r.Get("/", h.SessionsGet)
							r.Delete("/", h.SessionsTerminate)
						})
						r.Get("/", h.AccountGet)
						r.Post("/", h.Logout)
					})

					r.With(roleGrant).Route("/{account_id}", func(r chi.Router) {
						r.Route("/sessions", func(r chi.Router) {
							r.Get("/{session_id}", h.AdminSessionGet)
							r.Get("/", h.AdminSessionsGet)
							r.Delete("/", h.AdminSessionsTerminate)
						})
						r.Get("/", h.AdminAccountGet)
						r.Patch("/{role}", h.AdminRoleUpdate)
					})
				})
			})

			r.Post("/test/login", h.LoginSimple)
		})
	})

	server := httpkit.StartServer(ctx, svc.Config.Server.Port, r, svc.Log)

	<-ctx.Done()
	httpkit.StopServer(context.Background(), server, svc.Log)
}
