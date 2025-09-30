package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/errx"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (s *Service) DeleteUserSession(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid user id: %s", chi.URLParam(r, "user_id"))
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid user id: %s", chi.URLParam(r, "user_id")),
		})...)

		return
	}

	sessionID, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid session id: %s", chi.URLParam(r, "session_id"))
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid session id: %s", chi.URLParam(r, "session_id")),
		})...)

		return
	}

	if err = s.domain.Session.DeleteOneForUser(r.Context(), userID, sessionID); err != nil {
		s.log.WithError(err).Errorf("failed to delete user session")
		switch {
		case errors.Is(err, errx.ErrorSessionNotFound):
			ape.RenderErr(w, problems.NotFound("session for user not found"))
		case errors.Is(err, errx.ErrorNoPermissions):
			ape.RenderErr(w, problems.Forbidden("no permissions to delete session"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
