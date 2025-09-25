package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/api/rest/requests"
	"github.com/chains-lab/sso-svc/internal/api/rest/responses"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
)

func (s Service) Login(w http.ResponseWriter, r *http.Request) {
	req, err := requests.Login(r)
	if err != nil {
		s.log.WithError(err).Error("failed to decode login request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	token, err := s.domain.session.Login(r.Context(), req.Data.Attributes.Email, req.Data.Attributes.Password)
	if err != nil {
		s.log.WithError(err).Errorf("failed to login user")
		switch {
		case errors.Is(err, errx.ErrorUserNotFound):
			ape.RenderErr(w, problems.NotFound("user with this email not found"))
		case errors.Is(err, errx.ErrorInitiatorIsBlocked):
			ape.RenderErr(w, problems.Forbidden("user is blocked"))
		case errors.Is(err, errx.ErrorInvalidLogin):
			ape.RenderErr(w, problems.Unauthorized("invalid login or password"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	s.log.Infof("user %s logged in successfully", req.Data.Attributes.Email)

	ape.Render(w, http.StatusOK, responses.TokensPair(token))
}
