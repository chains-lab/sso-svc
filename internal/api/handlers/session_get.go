package handlers

import (
	"net/http"

	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/sso-oauth/internal/api/responses"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) SessionGet(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountData(r.Context())
	if err != nil {
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	session, err := h.app.GetSession(r.Context(), data.SessionID)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.Session(session))
}
