package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/api/rest/meta"
	"github.com/chains-lab/sso-svc/internal/api/rest/requests"
	"github.com/chains-lab/sso-svc/internal/api/rest/responses"
	"github.com/chains-lab/sso-svc/internal/errx"
)

func (s Service) RefreshToken(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.Log(r).WithError(err).Error("failed to get user from context")

		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))
		return
	}

	req, err := requests.RefreshSession(r)
	if err != nil {
		s.Log(r).WithError(err).Error("failed to parse refresh session request")

		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	tokensPair, err := s.app.RefreshSessionToken(r.Context(), initiator.UserID, initiator.SessionID, "TODO", "TODO", req.Data.Attributes.RefreshToken)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to refresh session token")
		switch {
		case errors.Is(err, errx.ErrorUnauthenticated):
			ape.RenderErr(w, problems.Unauthorized("failed to refresh session token"))
		case errors.Is(err, errx.ErrorInitiatorIsBlocked):
			ape.RenderErr(w, problems.Forbidden("user is blocked"))
		case errors.Is(err, errx.ErrorSessionNotFound):
			ape.RenderErr(w, problems.Unauthorized("session not found"))
		case errors.Is(err, errx.ErrorSessionClientMismatch):
			ape.RenderErr(w, problems.Unauthorized("session client mismatch"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.TokensPair(tokensPair))
}
