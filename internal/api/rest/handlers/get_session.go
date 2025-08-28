package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/api/rest/meta"
	"github.com/chains-lab/sso-svc/internal/api/rest/responses"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s Service) GetSession(w http.ResponseWriter, r *http.Request) {
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

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid user id: %s", chi.URLParam(r, "user_id"))

		ape.RenderErr(w, problems.InvalidParameter("user_id", err))
		return
	}

	session, err := s.app.AdminGetUserSession(r.Context(), initiator.UserID, userID, sessionID)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to get user session")

		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.UserSession(session))
}
