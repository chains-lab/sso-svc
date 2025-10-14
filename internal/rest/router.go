package rest

import (
	"context"
	"net/http"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/restkit/roles"
	"github.com/chains-lab/sso-svc/internal"
	"github.com/chains-lab/sso-svc/internal/rest/meta"
	"github.com/go-chi/chi/v5"
)

type Handlers interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	GoogleLogin(w http.ResponseWriter, r *http.Request)
	GoogleLoginCallback(w http.ResponseWriter, r *http.Request)
	RefreshSession(w http.ResponseWriter, r *http.Request)
	GetMyUser(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)
	GetMySessions(w http.ResponseWriter, r *http.Request)
	DeleteMySessions(w http.ResponseWriter, r *http.Request)
	GetMySession(w http.ResponseWriter, r *http.Request)
	DeleteMySession(w http.ResponseWriter, r *http.Request)
	RegisterAdmin(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	GetUserSessions(w http.ResponseWriter, r *http.Request)
	DeleteUserSessions(w http.ResponseWriter, r *http.Request)
	GetSession(w http.ResponseWriter, r *http.Request)
	DeleteUserSession(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	ServiceGrant(serviceName, skService string) func(http.Handler) http.Handler
	Auth(userCtxKey interface{}, skUser string) func(http.Handler) http.Handler
	RoleGrant(userCtxKey interface{}, allowedRoles map[string]bool) func(http.Handler) http.Handler
}

func Run(ctx context.Context, cfg internal.Config, log logium.Logger, m Middlewares, h Handlers) {
	svcAuth := m.ServiceGrant(cfg.Service.Name, cfg.JWT.Service.SecretKey)
	userAuth := m.Auth(meta.UserCtxKey, cfg.JWT.User.AccessToken.SecretKey)
	sysadmin := m.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Admin: true,
	})

	r := chi.NewRouter()

	r.Route("/sso-svc", func(r chi.Router) {
		r.Use(svcAuth)

		r.Route("/v1", func(r chi.Router) {
			r.Post("/register", h.RegisterUser)

			r.Route("/login", func(r chi.Router) {
				r.Post("/", h.Login)

				r.Route("/google", func(r chi.Router) {
					r.Post("/", h.GoogleLogin)
					r.Post("/callback", h.GoogleLoginCallback)
				})
			})

			r.Post("/refresh", h.RefreshSession)

			r.With(userAuth).Route("/me", func(r chi.Router) {
				r.With(userAuth).Get("/", h.GetMyUser)
				r.With(userAuth).Post("/logout", h.Logout)
				r.With(userAuth).Post("/password", h.ResetPassword)

				r.With(userAuth).Route("/sessions", func(r chi.Router) {
					r.Get("/", h.GetMySessions)
					r.Delete("/", h.DeleteMySessions)

					r.Route("/{session_id}", func(r chi.Router) {
						r.Get("/", h.GetMySession)
						r.Delete("/", h.DeleteMySession)
					})
				})
			})

			r.Route("/admin", func(r chi.Router) {
				r.Use(userAuth)
				r.Use(sysadmin)

				r.Post("/", h.RegisterAdmin)

				r.Route("/{user_id}", func(r chi.Router) {
					r.Get("/", h.GetUser)

					r.Route("/sessions", func(r chi.Router) {
						r.Get("/", h.GetUserSessions)
						r.Delete("/", h.DeleteUserSessions)

						r.Route("/{session_id}", func(r chi.Router) {
							r.Get("/", h.GetSession)
							r.Delete("/", h.DeleteUserSession)
						})
					})
				})
			})
		})
	})

	log.Infof("starting REST service on %s", cfg.Rest.Port)

	<-ctx.Done()

	log.Info("shutting down REST service")
}
