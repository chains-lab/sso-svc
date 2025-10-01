package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/chains-lab/sso-svc/internal/rest/meta"
)

func (s *Service) Logout(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	err = s.domain.Session.DeleteOneForUser(r.Context(), initiator.ID, initiator.SessionID)
	if err != nil {
		s.log.WithError(err).Errorf("failed to logout user")
		switch {
		case errors.Is(err, errx.ErrorSessionNotFound):
			ape.RenderErr(w, problems.Unauthorized("current session not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
