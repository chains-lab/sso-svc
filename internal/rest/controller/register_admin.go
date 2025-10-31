package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/chains-lab/sso-svc/internal/rest/meta"
	"github.com/chains-lab/sso-svc/internal/rest/requests"
	"github.com/chains-lab/sso-svc/internal/rest/responses"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (s *Service) RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.RegisterAdmin(r)
	if err != nil {
		s.log.WithError(err).Error("failed to decode register admin request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	u, err := s.domain.Auth.RegisterAdmin(r.Context(), initiator.ID,
		req.Data.Attributes.Email,
		req.Data.Attributes.Password,
		req.Data.Attributes.Role,
	)
	if err != nil {
		s.log.WithError(err).Errorf("failed to register by admin")
		switch {
		case errors.Is(err, errx.ErrorInitiatorIsBlocked):
			ape.RenderErr(w, problems.Forbidden("initiator is blocked"))
		case errors.Is(err, errx.ErrorUserAlreadyExists):
			ape.RenderErr(w, problems.Conflict("user with this email already exists"))
		case errors.Is(err, errx.ErrorRoleNotSupported):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"data/attributes/role": err,
			})...)
		case errors.Is(err, errx.ErrorPasswordIsInappropriate):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"data/attributes/password": err,
			})...)
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	s.log.Infof("admin %s registered successfully by user %s", u.ID, initiator.ID)

	ape.Render(w, http.StatusCreated, responses.User(u))
}
