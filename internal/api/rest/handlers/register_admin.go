package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/api/rest/meta"
	"github.com/chains-lab/sso-svc/internal/api/rest/requests"
	"github.com/chains-lab/sso-svc/internal/api/rest/responses"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/errx"
)

func (s Handler) RegisterAdmin(w http.ResponseWriter, r *http.Request) {
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

	err = roles.ParseRole(req.Data.Attributes.Role)
	if err != nil {
		ape.RenderErr(w, problems.InvalidParameter("data/attributes/role", err))

		return
	}
	if req.Data.Attributes.Role == roles.User {
		ape.RenderErr(w, problems.InvalidParameter(
			"data/attributes/role",
			fmt.Errorf("cannot register user with role 'user'"),
		))

		return
	}

	user, err := s.app.RegisterAdmin(r.Context(), initiator.UserID, initiator.SessionID, app.RegisterAdminParams{
		Email:    req.Data.Attributes.Email,
		Password: req.Data.Attributes.Password,
		Role:     req.Data.Attributes.Role,
	})
	if err != nil {
		s.log.WithError(err).Errorf("failed to register admin")
		switch {
		case errors.Is(err, errx.ErrorNoPermissions):
			ape.RenderErr(w, problems.Forbidden("no permissions to register admin"))
		case errors.Is(err, errx.ErrorUnauthenticated):
			ape.RenderErr(w, problems.Unauthorized("failed to register admin"))
		case errors.Is(err, errx.ErrorInitiatorIsBlocked):
			ape.RenderErr(w, problems.Forbidden("initiator is blocked"))
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

	s.log.Infof("admin %s registered successfully by user %s", user.ID, initiator.UserID)

	ape.Render(w, http.StatusCreated, responses.User(user))
}
