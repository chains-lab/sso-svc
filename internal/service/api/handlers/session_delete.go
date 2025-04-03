package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/sso-oauth/internal/service/api/responses"
	"github.com/hs-zavet/tokens"
)

func SessionDelete(w http.ResponseWriter, r *http.Request) {
	accountID, sessionID, _, _, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		Log(r).Warnf("Unauthorized session delete attempt: %v", err)
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	sessionForDeleteId, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if sessionID.String() == sessionForDeleteId.String() {
		Log(r).Debugf("Sessions can't be current")
		httpkit.RenderErr(w, problems.BadRequest(errors.New("session can't be current"))...)
		return
	}

	err = Domain(r).SessionDelete(r.Context(), sessionForDeleteId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	sessions, err := Domain(r).SessionsListByAccount(r.Context(), *accountID)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.SessionCollection(sessions))
}
