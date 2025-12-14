package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/domain"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/chains-lab/sso-svc/internal/rest/meta"
	"github.com/chains-lab/sso-svc/internal/rest/requests"
	"github.com/chains-lab/sso-svc/internal/rest/responses"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (s *Service) RegistrationAdmin(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.AccountData(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.RegistrationAdmin(r)
	if err != nil {
		s.log.WithError(err).Error("failed to decode register admin request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	u, err := s.domain.Registration(r.Context(), domain.RegistrationParams{
		Username: req.Data.Attributes.Username,
		Email:    req.Data.Attributes.Email,
		Password: req.Data.Attributes.Password,
		Role:     req.Data.Attributes.Role,
	})
	if err != nil {
		s.log.WithError(err).Errorf("failed to register by admin")
		switch {
		case errors.Is(err, errx.ErrorInitiatorNotFound):
			ape.RenderErr(w, problems.Unauthorized("failed to register admin user not found"))
		case errors.Is(err, errx.ErrorInitiatorIsNotActive):
			ape.RenderErr(w, problems.Forbidden("initiator account is not active"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("only admins can register new admin accounts"))
		case errors.Is(err, errx.ErrorEmailAlreadyExist):
			ape.RenderErr(w, problems.Conflict("user with this email already exists"))
		case errors.Is(err, errx.ErrorUsernameAlreadyTaken):
			ape.RenderErr(w, problems.Conflict("user with this username already exists"))
		case errors.Is(err, errx.ErrorUsernameIsNotAllowed):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"repo/attributes/username": err,
			})...)
		case errors.Is(err, errx.ErrorPasswordIsNotAllowed):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"repo/attributes/password": err,
			})...)
		case errors.Is(err, errx.ErrorRoleNotSupported):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"repo/attributes/role": err,
			})...)
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	s.log.Infof("admin %s registered successfully by user %s", u.ID, initiator)

	ape.Render(w, http.StatusCreated, responses.Account(u))
}
