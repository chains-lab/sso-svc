package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/chains-lab/sso-svc/internal/rest/requests"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (s *Service) RegisterUser(w http.ResponseWriter, r *http.Request) {
	req, err := requests.RegisterUser(r)
	if err != nil {
		s.log.WithError(err).Error("failed to decode register admin request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	_, err = s.domain.Auth.Register(r.Context(),
		req.Data.Attributes.Email,
		req.Data.Attributes.Password,
		roles.User,
	)
	if err != nil {
		s.log.WithError(err).Errorf("failed to register admin")
		switch {
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

	s.log.Infof("user %s registered successfully", req.Data.Attributes.Email)

	w.WriteHeader(http.StatusCreated)
}
