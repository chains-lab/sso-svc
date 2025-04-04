package handlers

import (
	"net/http"

	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/sso-oauth/internal/api/responses"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) SessionsGet(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountData(r.Context())
	if err != nil {
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	sessions, err := h.app.GetSessions(r.Context(), data.AccountID)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.SessionCollection(sessions))
}
