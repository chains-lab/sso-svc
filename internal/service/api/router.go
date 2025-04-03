package api

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/sso-oauth/internal/service"
	"github.com/hs-zavet/sso-oauth/internal/service/api/handlers"
	"github.com/hs-zavet/tokens"
	"github.com/hs-zavet/tokens/identity"
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

	r.Route("/re-news/sso", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/public", func(r chi.Router) {
				r.Post("/refresh", handlers.Refresh)
				r.Route("/google", func(r chi.Router) {
					r.Get("/login", handlers.GoogleLogin)
					r.Get("/callback", handlers.GoogleCallback)
				})

				r.Route("/account", func(r chi.Router) {
					r.Use(authMW)
					r.Route("/sessions", func(r chi.Router) {
						r.Route("/{session_id}", func(r chi.Router) {
							r.Get("/", handlers.SessionGet)
							r.Delete("/", handlers.SessionDelete)
						})
						r.Get("/", handlers.SessionsGet)
						r.Delete("/", handlers.SessionsTerminate)
					})
					r.Get("/", handlers.AccountGet)
					r.Delete("/", handlers.Logout)
				})
			})

			r.Route("/private", func(r chi.Router) {
				r.Route("/account", func(r chi.Router) {
					r.Use(roleGrant)
					r.Route("/{account_id}", func(r chi.Router) {
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
