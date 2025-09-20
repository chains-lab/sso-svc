package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/api/rest/meta"
	"github.com/chains-lab/sso-svc/internal/api/rest/responses"
	"github.com/chains-lab/sso-svc/internal/errx"
)

func (s Handler) GetOwnUser(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	user, err := s.app.GetUserByID(r.Context(), initiator.UserID)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get user by id: %s", initiator.UserID)
		switch {
		case errors.Is(err, errx.ErrorUserNotFound):
			ape.RenderErr(w, problems.NotFound("user not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.User(user))
}
