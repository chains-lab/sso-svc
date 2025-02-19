package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/tokens"
)

func (a *App) SessionsTerminate(w http.ResponseWriter, r *http.Request) {
	accountID, sessionID, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		a.Log.Warnf("Unauthorized session terminate attempt: %v", err)
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	err = a.Domain.SessionsTerminate(r.Context(), *accountID, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		a.Log.Errorf("Failed to terminate session: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, http.StatusOK)
}
