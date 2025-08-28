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
	"github.com/chains-lab/sso-svc/internal/errx"
)

func (s Service) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.Log(r).WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	if initiator.Role != roles.User {
		ape.RenderErr(w, problems.Forbidden("only users can update their password"))

		return
	}

	req, err := requests.UpdatePassword(r)
	if err != nil {
		s.Log(r).WithError(err).Error("failed to decode update password request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	if req.Data.Attributes.NewPassword != req.Data.Attributes.ConfirmPassword {
		ape.RenderErr(w,
			problems.InvalidParameter(
				"data/attributes/new_password",
				fmt.Errorf("passwords and confirm do not match"),
			),
			problems.InvalidParameter(
				"data/attributes/confirm_password",
				fmt.Errorf("passwords and confirm do not match"),
			),
		)

		return
	}

	if req.Data.Attributes.OldPassword == req.Data.Attributes.NewPassword {
		ape.RenderErr(w,
			problems.InvalidParameter(
				"data/attributes/new_password",
				fmt.Errorf("new password must be different from the current password"),
			),
		)

		return
	}

	err = s.app.UpdatePassword(r.Context(), initiator.UserID, initiator.SessionID, req.Data.Attributes.OldPassword, req.Data.Attributes.NewPassword)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to update password")
		switch {
		case errors.Is(err, errx.ErrorUnauthenticated):
			ape.RenderErr(w, problems.Unauthorized("failed to update password"))
		case errors.Is(err, errx.ErrorInitiatorIsBlocked):
			ape.RenderErr(w, problems.Forbidden("user is blocked"))
		case errors.Is(err, errx.ErrorInvalidLogin):
			ape.RenderErr(w, problems.InvalidParameter(
				"data/attributes/old_password",
				fmt.Errorf("current password is incorrect"),
			))
		case errors.Is(err, errx.ErrorPasswordIsInappropriate):
			ape.RenderErr(w, problems.InvalidParameter(
				"data/attributes/new_password",
				fmt.Errorf("new password does not meet complexity requirements"),
			))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
