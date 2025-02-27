package handlers

import (
	"net/http"

	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/api/responses"
	"github.com/recovery-flow/tokens"
)

func (h *Handler) SessionsGet(w http.ResponseWriter, r *http.Request) {
	accountID, _, _, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		h.Log.WithError(err).Warn("Unauthorized session list attempt")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	sessions, err := h.Domain.SessionsListByAccount(r.Context(), *accountID)
	if err != nil {
		h.Log.WithError(err).Error("Failed to list sessions")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.SessionCollection(sessions))
}
