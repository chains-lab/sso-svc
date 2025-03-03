package handlers

import (
	"net/http"

	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/api/responses"
	"github.com/recovery-flow/tokens"
)

func SessionGet(w http.ResponseWriter, r *http.Request) {
	accountID, sessionID, _, _, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		Log(r).WithError(err).Debug("Failed to get account data")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	session, err := Domain(r).SessionGetForAccount(r.Context(), *sessionID, *accountID)
	if err != nil {
		Log(r).WithError(err).Debug("Failed to get session")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	if session.AccountID != *accountID {
		Log(r).Errorf("Session doesn't belong to account")
		httpkit.RenderErr(w, problems.Forbidden("Session doesn't belong to account"))
		return
	}

	httpkit.Render(w, responses.Session(*session))
}
