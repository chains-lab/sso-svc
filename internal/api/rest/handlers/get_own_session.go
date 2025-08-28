package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/api/rest/meta"
	"github.com/chains-lab/sso-svc/internal/api/rest/responses"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s Service) GetOwnSession(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.Log(r).WithError(err).Error("failed to get user from context")

		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))
		return
	}

	sessionId, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid session id: %s", chi.URLParam(r, "session_id"))

		ape.RenderErr(w, problems.InvalidParameter("session_id", err))
		return
	}

	session, err := s.app.GetOwnSession(r.Context(), initiator.UserID, sessionId)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to get own session")

		switch {
		case errors.Is(err, errx.ErrorSessionNotFound):
			ape.RenderErr(w, problems.NotFound("session not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.UserSession(session))
}
