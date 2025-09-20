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

type Handlers interface {
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
	DeleteSessions(w http.ResponseWriter, r *http.Request)
	GetSession(w http.ResponseWriter, r *http.Request)
	DeleteSession(w http.ResponseWriter, r *http.Request)
}

func (a *Service) Run(ctx context.Context, h Handlers) {
	svcAuth := mdlv.ServiceAuthMdl(enum.SsoSVC, a.cfg.JWT.Service.SecretKey)
	userAuth := mdlv.AuthMdl(meta.UserCtxKey, a.cfg.JWT.User.AccessToken.SecretKey)
	adminGrant := mdlv.AccessGrant(meta.UserCtxKey, roles.Admin, roles.SuperUser)

	a.router.Route("/sso-svc/", func(r chi.Router) {
		r.Use(svcAuth)
		r.Route("/v1", func(r chi.Router) {
			r.Post("/register", h.RegisterUser)

			r.Route("/login", func(r chi.Router) {
				r.Post("/", h.Login)

				r.Route("/google", func(r chi.Router) {
					r.Post("/", h.GoogleLogin)
					r.Post("/callback", h.GoogleCallback)
				})
			})

			r.Post("/refresh", h.RefreshToken)

			r.With(userAuth).Route("/own", func(r chi.Router) {
				r.With(userAuth).Get("/", h.GetOwnUser)
				r.With(userAuth).Post("/logout", h.Logout)
				r.With(userAuth).Post("/password", h.UpdatePassword)

				r.With(userAuth).Route("/sessions", func(r chi.Router) {
					r.Get("/", h.SelectOwnSessions)
					r.Delete("/", h.DeleteOwnSessions)

					r.Route("/{session_id}", func(r chi.Router) {
						r.Get("/", h.GetOwnSession)
						r.Delete("/", h.DeleteOwnSession)
					})
				})
			})

			r.Route("/admin", func(r chi.Router) {
				r.Use(userAuth)
				r.Use(adminGrant)

				r.Post("/", h.RegisterAdmin)

				r.Route("/{user_id}", func(r chi.Router) {
					r.Get("/", h.GetUser)

					r.Route("/sessions", func(r chi.Router) {
						r.Get("/", h.SelectUserSessions)
						r.Delete("/", h.DeleteSessions)

						r.Route("/{session_id}", func(r chi.Router) {
							r.Get("/", h.GetSession)
							r.Delete("/", h.DeleteSession)
						})
					})
				})
			})
		})
	})

	a.Start(ctx)

	<-ctx.Done()
	a.Stop(ctx)
}
