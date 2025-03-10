package handlers

import (
	"net/http"

	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/tokens"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	accountID, sessionID, _, _, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		Log(r).WithError(err).Warn("Unauthorized logout attempt")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	err = Domain(r).SessionDelete(r.Context(), *sessionID)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	Log(r).Infof("Account %s logged out", accountID)
	httpkit.Render(w, http.StatusNoContent)
}
