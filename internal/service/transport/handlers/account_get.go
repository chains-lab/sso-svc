package handlers

import (
	"net/http"

	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/responses"
	"github.com/recovery-flow/tokens"
)

func (h *Handlers) AccountGet(w http.ResponseWriter, r *http.Request) {
	accountID, _, _, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		h.Log.WithError(err).Error("Failed to retrieve account data")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	account, err := h.Domain.AccountGet(r.Context(), *accountID)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}
	httpkit.Render(w, responses.Account(*account))
}
