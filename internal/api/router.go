package api

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/hs-zavet/comtools/httpkit"
	handlers2 "github.com/hs-zavet/sso-oauth/internal/api/handlers"
	"github.com/hs-zavet/sso-oauth/internal/service"
	"github.com/hs-zavet/tokens"
	"github.com/hs-zavet/tokens/identity"
)

func Run(ctx context.Context, svc *service.Service) {
	r := chi.NewRouter()

	r.Use(
		httpkit.CtxMiddleWare(
			handlers2.CtxLog(svc.Log),
			handlers2.CtxDomain(svc.Domain),
			handlers2.CtxConfig(svc.Config),
			handlers2.CtxGoogleOauth(svc.Config),
		),
	)

	authMW := tokens.AuthMdl(svc.Config.JWT.AccessToken.SecretKey)
	roleGrant := tokens.IdentityMdl(svc.Config.JWT.AccessToken.SecretKey, identity.Admin, identity.SuperUser)

	r.Route("/re-news/sso", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/public", func(r chi.Router) {
				r.Post("/refresh", handlers2.Refresh)
				r.Route("/google", func(r chi.Router) {
					r.Get("/login", handlers2.GoogleLogin)
					r.Get("/callback", handlers2.GoogleCallback)
				})

				r.Route("/account", func(r chi.Router) {
					r.Use(authMW)
					r.Route("/sessions", func(r chi.Router) {
						r.Route("/{session_id}", func(r chi.Router) {
							r.Get("/", handlers2.SessionGet)
							r.Delete("/", handlers2.SessionDelete)
						})
						r.Get("/", handlers2.SessionsGet)
						r.Delete("/", handlers2.SessionsTerminate)
					})
					r.Get("/", handlers2.AccountGet)
					r.Delete("/", handlers2.Logout)
				})
			})

			r.Route("/private", func(r chi.Router) {
				r.Route("/account", func(r chi.Router) {
					r.Use(roleGrant)
					r.Route("/{account_id}", func(r chi.Router) {
						r.Route("/sessions", func(r chi.Router) {
							r.Get("/{session_id}", handlers2.AdminSessionGet)
							r.Get("/", handlers2.AdminSessionsGet)
							r.Delete("/", handlers2.AdminSessionsTerminate)
						})
						r.Get("/", handlers2.AdminAccountGet)
						r.Patch("/{role}", handlers2.AdminRoleUpdate)
					})
				})
			})

			r.Post("/test/login", handlers2.LoginSimple)
		})
	})

	server := httpkit.StartServer(ctx, svc.Config.Server.Port, r, svc.Log)

	<-ctx.Done()
	httpkit.StopServer(context.Background(), server, svc.Log)
}
