package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/api/rest/meta"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s Service) DeleteOwnSession(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.Log(r).WithError(err).Error("failed to get user from context")

		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))
		return
	}

	sessionID, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid session id: %s", chi.URLParam(r, "session_id"))

		ape.RenderErr(w, problems.InvalidParameter("session_id", err))
		return
	}

	if err := s.app.DeleteOwnSession(r.Context(), initiator.UserID, sessionID); err != nil {
		s.Log(r).WithError(err).Errorf("failed to delete own session")

		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusNoContent, nil)
}
