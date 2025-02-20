package handlers

import (
	"net/http"

	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/responses"
	"github.com/recovery-flow/tokens"
)

func (a *App) AccountGet(w http.ResponseWriter, r *http.Request) {
	accountID, _, _, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		a.Log.Warnf("Unauthorized account get attempt: %v", err)
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	user, err := a.Domain.AccountGet(r.Context(), *accountID)
	if err != nil {
		a.Log.Errorf("Failed to retrieve user: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}
	httpkit.Render(w, responses.Account(*user))
}
