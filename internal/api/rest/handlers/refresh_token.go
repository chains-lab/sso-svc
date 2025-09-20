package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/api/rest/requests"
	"github.com/chains-lab/sso-svc/internal/api/rest/responses"
	"github.com/chains-lab/sso-svc/internal/errx"
)

func (s Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	req, err := requests.RefreshSession(r)
	if err != nil {
		s.log.WithError(err).Error("failed to parse refresh session request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	tokensPair, err := s.app.RefreshSession(r.Context(), req.Data.Attributes.RefreshToken)
	if err != nil {
		s.log.WithError(err).Errorf("failed to refresh session token")
		switch {
		case errors.Is(err, errx.ErrorUnauthenticated):
			ape.RenderErr(w, problems.Unauthorized("failed to refresh session token"))
		case errors.Is(err, errx.ErrorInitiatorIsBlocked):
			ape.RenderErr(w, problems.Forbidden("user is blocked"))
		case errors.Is(err, errx.ErrorSessionNotFound):
			ape.RenderErr(w, problems.Unauthorized("session not found"))
		case errors.Is(err, errx.ErrorSessionTokenMismatch):
			ape.RenderErr(w, problems.Unauthorized("refresh session token mismatch"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.TokensPair(tokensPair))
}
