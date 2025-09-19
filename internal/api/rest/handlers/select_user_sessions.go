package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/api/rest/meta"
	"github.com/chains-lab/sso-svc/internal/api/rest/responses"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s Service) SelectUserSessions(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.Log(r).WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	pagReq, sort := pagi.GetPagination(r)

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid user id: %s", chi.URLParam(r, "user_id"))
		ape.RenderErr(w, problems.InvalidParameter("user_id", err))

		return
	}

	sessions, pag, err := s.app.AdminListUserSessions(r.Context(), initiator.UserID, initiator.SessionID, userID, pagReq, sort)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to select own sessions")
		switch {
		case errors.Is(err, errx.ErrorUnauthenticated):
			ape.RenderErr(w, problems.Unauthorized("failed to select user sessions"))
		case errors.Is(err, errx.ErrorInitiatorIsBlocked):
			ape.RenderErr(w, problems.Forbidden("user is blocked"))
		case errors.Is(err, errx.ErrorNoPermissions):
			ape.RenderErr(w, problems.Forbidden("no permissions to select user sessions"))
		case errors.Is(err, errx.ErrorUserNotFound):
			ape.RenderErr(w, problems.NotFound("user not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.UserSessionsCollection(sessions, pag))
}
