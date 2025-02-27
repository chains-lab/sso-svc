package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/tokens"
)

func (h *Handler) SessionsTerminate(w http.ResponseWriter, r *http.Request) {
	accountID, sessionID, _, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		h.Log.WithError(err).Warn("Unauthorized session terminate attempt")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	err = h.Domain.SessionsTerminate(r.Context(), *accountID, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, http.StatusOK)
}
