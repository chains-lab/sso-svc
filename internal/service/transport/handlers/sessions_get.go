package handlers

import (
	"net/http"

	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/responses"
	"github.com/recovery-flow/tokens"
)

func (a *App) SessionsGet(w http.ResponseWriter, r *http.Request) {
	accountID, _, _, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		a.Log.Warnf("Unauthorized session get attempt: %v", err)
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	sessions, err := a.Domain.SessionsListByUser(r.Context(), *accountID)
	if err != nil {
		a.Log.Errorf("Failed to retrieve user sessions: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.SessionCollection(sessions))
}
