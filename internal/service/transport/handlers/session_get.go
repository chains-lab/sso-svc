package handlers

import (
	"net/http"

	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/responses"
	"github.com/recovery-flow/tokens"
)

func (a *App) SessionGet(w http.ResponseWriter, r *http.Request) {
	accountID, sessionID, _, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		a.Log.Warnf("Unauthorized session get attempt: %v", err)
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	session, err := a.Domain.SessionGetForUser(r.Context(), *sessionID, *accountID)
	if err != nil {
		a.Log.Errorf("Failed to retrieve user session: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	if session.UserID != *accountID {
		a.Log.Debugf("Session doesn't belong to user")
		httpkit.RenderErr(w, problems.Forbidden("Session doesn't belong to user"))
		return
	}

	httpkit.Render(w, responses.Session(*session))
}
