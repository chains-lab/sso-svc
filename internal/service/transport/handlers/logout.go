package handlers

import (
	"net/http"

	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/core/tools"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	sessionID, userID, err := tools.GetSessionAndUserID(r.Context())
	if err != nil {
		Log.Warnf("Unauthorized logout attempt: %v", err)
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	err = TokenManager.AddToBlackList(r.Context(), sessionID.String(), userID.String())
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	err = Domain.Session.Delete(r.Context(), sessionID)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, http.StatusOK)
}
