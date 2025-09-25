package rest

import (
	"context"
	"net/http"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/mdlv"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/api/rest/meta"
	"github.com/go-chi/chi/v5"
)

type Controller interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	GoogleLogin(w http.ResponseWriter, r *http.Request)
	GoogleCallback(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	GetOwnUser(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	UpdatePassword(w http.ResponseWriter, r *http.Request)
	SelectOwnSessions(w http.ResponseWriter, r *http.Request)
	DeleteOwnSessions(w http.ResponseWriter, r *http.Request)
	GetOwnSession(w http.ResponseWriter, r *http.Request)
	DeleteOwnSession(w http.ResponseWriter, r *http.Request)
	RegisterAdmin(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	SelectUserSessions(w http.ResponseWriter, r *http.Request)
	DeleteUserSessions(w http.ResponseWriter, r *http.Request)
	GetSession(w http.ResponseWriter, r *http.Request)
	DeleteUserSession(w http.ResponseWriter, r *http.Request)
}

func (s *Service) Run(ctx context.Context, c Controller) {
	svcAuth := mdlv.ServiceGrant(enum.SsoSVC, s.cfg.JWT.Service.SecretKey)
	userAuth := mdlv.Auth(meta.UserCtxKey, s.cfg.JWT.User.AccessToken.SecretKey)
	sysadmin := mdlv.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Admin: true,
	})

	s.router.Route("/sso-svc/", func(r chi.Router) {
		r.Use(svcAuth)
		r.Route("/v1", func(r chi.Router) {
			r.Post("/register", c.RegisterUser)

			r.Route("/login", func(r chi.Router) {
				r.Post("/", c.Login)

				r.Route("/google", func(r chi.Router) {
					r.Post("/", c.GoogleLogin)
					r.Post("/callback", c.GoogleCallback)
				})
			})

			r.Post("/refresh", c.RefreshToken)

			r.With(userAuth).Route("/own", func(r chi.Router) {
				r.With(userAuth).Get("/", c.GetOwnUser)
				r.With(userAuth).Post("/logout", c.Logout)
				r.With(userAuth).Post("/password", c.UpdatePassword)

				r.With(userAuth).Route("/sessions", func(r chi.Router) {
					r.Get("/", c.SelectOwnSessions)
					r.Delete("/", c.DeleteOwnSessions)

					r.Route("/{session_id}", func(r chi.Router) {
						r.Get("/", c.GetOwnSession)
						r.Delete("/", c.DeleteOwnSession)
					})
				})
			})

			r.Route("/admin", func(r chi.Router) {
				r.Use(userAuth)
				r.Use(sysadmin)

				r.Post("/", c.RegisterAdmin)

				r.Route("/{user_id}", func(r chi.Router) {
					r.Get("/", c.GetUser)

					r.Route("/sessions", func(r chi.Router) {
						r.Get("/", c.SelectUserSessions)
						r.Delete("/", c.DeleteUserSessions)

						r.Route("/{session_id}", func(r chi.Router) {
							r.Get("/", c.GetSession)
							r.Delete("/", c.DeleteUserSession)
						})
					})
				})
			})
		})
	})

	s.Start(ctx)

	<-ctx.Done()
	s.Stop(ctx)
}
