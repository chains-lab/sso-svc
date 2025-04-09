package handlers

import (
	"errors"
	"net/http"

	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/sso-oauth/internal/app/ape"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	err = h.app.Logout(r.Context(), data.SessionID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrSessionNotFound):
			h.log.WithError(err).Errorf("session not found session id: %s", data.SessionID)
			httpkit.RenderErr(w, problems.NotFound())
		case errors.Is(err, ape.ErrAccountNotFound):
			h.log.WithError(err).Errorf("account not found account id: %s", data.AccountID)
			httpkit.RenderErr(w, problems.NotFound("account not found"))
		default:
			h.log.WithError(err).Errorf("error logging out session id: %s", data.SessionID)
			httpkit.RenderErr(w, problems.InternalError())
		}
		return
	}

	httpkit.Render(w, http.StatusNoContent)
}
