package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/tokens"
)

func SessionsTerminate(w http.ResponseWriter, r *http.Request) {
	accountID, _, _, _, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		Log(r).WithError(err).Warn("Unauthorized session terminate attempt")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	err = Domain(r).SessionsTerminate(r.Context(), *accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, http.StatusNoContent)
}
