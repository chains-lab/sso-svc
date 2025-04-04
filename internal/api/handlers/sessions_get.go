package handlers

import (
	"net/http"

	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/sso-oauth/internal/api/responses"
	"github.com/hs-zavet/tokens"
)

func SessionsGet(w http.ResponseWriter, r *http.Request) {
	accountID, _, _, _, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		Log(r).WithError(err).Warn("Unauthorized session list attempt")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	sessions, err := Domain(r).SessionsListByAccount(r.Context(), *accountID)
	if err != nil {
		Log(r).WithError(err).Error("Failed to list sessions")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.SessionCollection(sessions))
}
