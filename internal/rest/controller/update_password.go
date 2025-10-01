package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/chains-lab/sso-svc/internal/rest/meta"
	"github.com/chains-lab/sso-svc/internal/rest/requests"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (s *Service) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.UpdatePassword(r)
	if err != nil {
		s.log.WithError(err).Error("failed to decode update password request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	err = s.domain.Auth.UpdatePassword(r.Context(), initiator.ID, req.Data.Attributes.OldPassword, req.Data.Attributes.NewPassword)
	if err != nil {
		s.log.WithError(err).Errorf("failed to update password")
		switch {
		case errors.Is(err, errx.ErrorUserNotFound):
			ape.RenderErr(w, problems.Unauthorized("failed to update password user not found"))
		case errors.Is(err, errx.ErrorInitiatorIsBlocked):
			ape.RenderErr(w, problems.Forbidden("initiator is blocked"))
		case errors.Is(err, errx.ErrorInvalidLogin):
			ape.RenderErr(w, problems.Forbidden("invalid credentials"))
		case errors.Is(err, errx.ErrorPasswordIsInappropriate):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"data/attributes/password": err,
			})...)
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
