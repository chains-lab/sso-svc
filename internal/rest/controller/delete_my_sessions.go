package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/chains-lab/sso-svc/internal/rest/meta"
)

func (s *Service) DeleteMySessions(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	if err = s.domain.Session.DeleteAllForUser(r.Context(), initiator.ID); err != nil {
		s.log.WithError(err).Errorf("failed to delete My sessions")
		switch {
		case errors.Is(err, errx.ErrorUnauthenticated):
			ape.RenderErr(w, problems.Unauthorized("failed to authenticate user"))
		case errors.Is(err, errx.ErrorUserNotFound):
			ape.RenderErr(w, problems.NotFound("user not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
