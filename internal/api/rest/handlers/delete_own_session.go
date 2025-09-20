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

func (s Handler) DeleteOwnSession(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	sessionID, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid session id: %s", chi.URLParam(r, "session_id"))
		ape.RenderErr(w, problems.InvalidParameter("session_id", err))

		return
	}

	if err := s.app.DeleteOwnSession(r.Context(), initiator.UserID, initiator.SessionID, sessionID); err != nil {
		s.log.WithError(err).Errorf("failed to delete own session")
		switch {
		case errors.Is(err, errx.ErrorUnauthenticated):
			ape.RenderErr(w, problems.Unauthorized("failed to authenticate user"))
		case errors.Is(err, errx.ErrorInitiatorIsBlocked):
			ape.RenderErr(w, problems.Forbidden("initiator is blocked"))
		case errors.Is(err, errx.ErrorSessionNotFound):
			ape.RenderErr(w, problems.NotFound("session not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
