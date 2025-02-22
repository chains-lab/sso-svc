package handlers

import (
	"net/http"

	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/api/responses"
	"github.com/recovery-flow/tokens"
)

func (h *Handlers) SessionGet(w http.ResponseWriter, r *http.Request) {
	accountID, sessionID, _, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		h.Log.WithError(err).Debug("Failed to get account data")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	session, err := h.Domain.SessionGetForAccount(r.Context(), *sessionID, *accountID)
	if err != nil {
		h.Log.WithError(err).Debug("Failed to get session")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	if session.AccountID != *accountID {
		h.Log.Errorf("Session doesn't belong to account")
		httpkit.RenderErr(w, problems.Forbidden("Session doesn't belong to account"))
		return
	}

	httpkit.Render(w, responses.Session(*session))
}
