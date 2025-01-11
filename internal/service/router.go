package service

import (
	"context"
	"time"

	"github.com/go-chi/chi"
	"github.com/recovery-flow/comtools/cifractx"
	"github.com/recovery-flow/comtools/httpkit"
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
	rateLimiter := httpkit.NewRateLimiter(15, 10*time.Second, 5*time.Minute)
	r.Route("/re-flow", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Use(rateLimiter.Middleware)
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
				r.Route("/user", func(r chi.Router) {
					r.Use(authMW)
					r.Route("/sessions", func(r chi.Router) {
						r.Route("/{session_id}", func(r chi.Router) {
							r.Get("/", handlers.GetSessions)
							r.Delete("/", handlers.DeleteSession)
						})
						r.Get("/", handlers.GetSessions)
						r.Delete("/terminate", handlers.TerminateSessions)
					})
					r.Post("/logout", handlers.Logout)
				})
			})
		})
	})

	server := httpkit.StartServer(ctx, service.Config.Server.Port, r, service.Logger)

	<-ctx.Done()
	httpkit.StopServer(context.Background(), server, service.Logger)
}
