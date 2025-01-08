package service

import (
	"context"
	"time"

	"github.com/cifra-city/comtools/cifractx"
	"github.com/cifra-city/comtools/httpkit"
	"github.com/cifra-city/sso-oauth/internal/config"
	"github.com/cifra-city/sso-oauth/internal/service/handlers"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

func Run(ctx context.Context) {
	r := chi.NewRouter()

	service, err := cifractx.GetValue[*config.Server](ctx, config.SERVER)
	if err != nil {
		logrus.Fatalf("failed to get server from context: %v", err)
	}

	r.Use(cifractx.MiddlewareWithContext(config.SERVER, service))
	authMW := service.TokenManager.AuthMiddleware(service.Config.JWT.AccessToken.SecretKey)
	rateLimiter := httpkit.NewRateLimiter(15, 10*time.Second, 5*time.Minute)
	r.Route("/cifra-sso", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Use(rateLimiter.Middleware)
			r.Route("/public", func(r chi.Router) {
				r.Post("/approve-operation", handlers.ApproveOperation)
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
						r.Get("/", handlers.GetSessions)
						r.Delete("/", handlers.DeleteSession)
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
