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

	r.Use(
		httpkit.CtxMiddleWare(
			handlers.CtxLog(svc.Log),
			handlers.CtxDomain(svc.Domain),
			handlers.CtxConfig(svc.Config),
			handlers.CtxGoogleOauth(svc.Config),
		),
	)

	authMW := tokens.AuthMdl(svc.Config.JWT.AccessToken.SecretKey)
	roleGrant := tokens.IdentityMdl(svc.Config.JWT.AccessToken.SecretKey, identity.Admin, identity.SuperUser)

	r.Route("/re-flow/sso", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/public", func(r chi.Router) {
				r.Post("/refresh", handlers.Refresh)
				r.Route("/google", func(r chi.Router) {
					r.Get("/login", handlers.GoogleLogin)
					r.Get("/callback", handlers.GoogleCallback)
				})
			})

			r.Route("/private", func(r chi.Router) {
				r.Route("/account", func(r chi.Router) {
					r.With(authMW).Route("/current", func(r chi.Router) {
						r.Route("/sessions", func(r chi.Router) {
							r.Route("/{session_id}", func(r chi.Router) {
								r.Get("/", handlers.SessionGet)
								r.Delete("/", handlers.SessionDelete)
							})
							r.Get("/", handlers.SessionsGet)
							r.Delete("/", handlers.SessionsTerminate)
						})
						r.Get("/", handlers.AccountGet)
						r.Post("/", handlers.Logout)
					})

					r.With(roleGrant).Route("/{account_id}", func(r chi.Router) {
						r.Route("/sessions", func(r chi.Router) {
							r.Get("/{session_id}", handlers.AdminSessionGet)
							r.Get("/", handlers.AdminSessionsGet)
							r.Delete("/", handlers.AdminSessionsTerminate)
						})
						r.Get("/", handlers.AdminAccountGet)
						r.Patch("/{role}", handlers.AdminRoleUpdate)
					})
				})
			})

			r.Post("/test/login", handlers.LoginSimple)
		})
	})

	server := httpkit.StartServer(ctx, svc.Config.Server.Port, r, svc.Log)

	<-ctx.Done()
	httpkit.StopServer(context.Background(), server, svc.Log)
}
