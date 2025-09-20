package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/api/rest/requests"
	"github.com/chains-lab/sso-svc/internal/errx"
)

func (s Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	req, err := requests.RegisterUser(r)
	if err != nil {
		s.log.WithError(err).Error("failed to decode register admin request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	if req.Data.Attributes.Password != req.Data.Attributes.ConfirmPassword {
		ape.RenderErr(w,
			problems.InvalidParameter(
				"data/attributes/confirm_password",
				fmt.Errorf("passwords and confirm do not match"),
			),
			problems.InvalidParameter(
				"data/attributes/password",
				fmt.Errorf("passwords and confirm do not match"),
			),
		)

		return
	}

	err = s.app.Register(r.Context(), req.Data.Attributes.Email, req.Data.Attributes.Password)
	if err != nil {
		s.log.WithError(err).Errorf("failed to register admin")
		switch {
		case errors.Is(err, errx.ErrorUserAlreadyExists):
			ape.RenderErr(w, problems.Conflict("user with this email already exists"))
		case errors.Is(err, errx.ErrorRoleNotSupported):
			ape.RenderErr(w, problems.InvalidParameter("data/attributes/role", err))
		case errors.Is(err, errx.ErrorPasswordIsInappropriate):
			ape.RenderErr(w, problems.InvalidParameter("data/attributes/password", err))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	s.log.Infof("user %s registered successfully", req.Data.Attributes.Email)

	w.WriteHeader(http.StatusCreated)
}
