package handlers

import (
	"net/http"

	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountData(r.Context())
	if err != nil {
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	err = h.app.DeleteSession(r.Context(), data.SessionID)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, http.StatusNoContent)
}
