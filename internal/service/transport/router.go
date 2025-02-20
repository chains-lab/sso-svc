package transport

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/handlers"
	"github.com/recovery-flow/tokens"
	"github.com/recovery-flow/tokens/identity"
)

func Run(ctx context.Context, app *handlers.App) {
	r := chi.NewRouter()

	authMW := tokens.AuthMdl(ctx, app.Config.JWT.AccessToken.SecretKey)
	adminGrant := tokens.IdentityMdl(ctx, app.Config.JWT.AccessToken.SecretKey, identity.Admin, identity.SuperUser)

	r.Route("/re-flow/sso", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/public", func(r chi.Router) {
				r.Post("/refresh", app.Refresh)

				r.Route("/oauth", func(r chi.Router) {
					r.Route("/google", func(r chi.Router) {
						r.Get("/login", app.GoogleLogin)
						r.Get("/callback", app.GoogleCallback)
					})
				})
			})

			r.Route("/private", func(r chi.Router) {
				r.Use(authMW)
				r.Route("/account", func(r chi.Router) {
					r.Get("/", app.AccountGet)

					r.Route("/sessions", func(r chi.Router) {
						r.Route("/{session_id}", func(r chi.Router) {
							r.Get("/", app.SessionGet)
							r.Delete("/", app.SessionDelete)
						})

						r.Get("/", app.SessionsGet)
						r.Delete("/", app.SessionsTerminate)
					})
					r.Post("/logout", app.Logout)
				})
			})

			r.Route("/admin", func(r chi.Router) {
				r.Use(adminGrant)
				r.Route("/{account_id}", func(r chi.Router) {
					r.Route("/sessions", func(r chi.Router) {
						r.Get("/", app.AdminSessionsGet)
						r.Delete("/", app.AdminSessionsTerminate)
					})

					r.Patch("/role/{role}", app.AdminRoleUpdate)
				})

				r.Route("/sessions", func(r chi.Router) {
					r.Route("/{session_id}", func(r chi.Router) {
						r.Get("/", app.AdminSessionsGet)
						r.Delete("/", app.AdminSessionsTerminate)
					})
				})
			})

			r.Post("/test/login", app.LoginSimple)
		})
	})

	server := httpkit.StartServer(ctx, app.Config.Server.Port, r, app.Log)

	<-ctx.Done()
	httpkit.StopServer(context.Background(), server, app.Log)
}
