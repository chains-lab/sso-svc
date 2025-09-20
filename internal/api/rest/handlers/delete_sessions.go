package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/api/rest/meta"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s Handler) DeleteSessions(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid user id: %s", chi.URLParam(r, "user_id"))
		ape.RenderErr(w, problems.InvalidParameter("user_id", err))

		return
	}

	if err := s.app.AdminDeleteUserSessions(r.Context(), initiator.UserID, initiator.SessionID, userID); err != nil {
		s.log.WithError(err).Errorf("failed to delete user sessions")
		switch {
		case errors.Is(err, errx.ErrorUnauthenticated):
			ape.RenderErr(w, problems.Unauthorized("unauthenticated"))
		case errors.Is(err, errx.ErrorInitiatorIsBlocked):
			ape.RenderErr(w, problems.Forbidden("initiator is blocked"))
		case errors.Is(err, errx.ErrorUserNotFound):
			ape.RenderErr(w, problems.NotFound("user not found"))
		case errors.Is(err, errx.ErrorNoPermissions):
			ape.RenderErr(w, problems.Forbidden("no permissions"))
		case errors.Is(err, errx.ErrorRoleNotSupported):
			ape.RenderErr(w, problems.InternalError()) //TODO
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
