package handlers

import (
	"net/http"

	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/tokens"
)

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	accountID, sessionID, _, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		h.Log.WithError(err).Warn("Unauthorized logout attempt")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	err = h.Domain.SessionDelete(r.Context(), *sessionID)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	h.Log.Infof("User %s logged out", accountID)
	httpkit.Render(w, http.StatusOK)
}
