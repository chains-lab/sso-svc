package handlers

import (
	"net/http"

	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/tokens"
)

func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	accountID, sessionID, _, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		a.Log.Warnf("Unauthorized logout attempt: %v", err)
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	err = a.Domain.SessionDelete(r.Context(), *sessionID)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	a.Log.Infof("User %s logged out", accountID)
	httpkit.Render(w, http.StatusOK)
}
