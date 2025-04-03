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
	"github.com/hs-zavet/tokens/identity"
)

func AdminSessionDelete(w http.ResponseWriter, r *http.Request) {
	initiatorID, initiatorSession, _, initiatorRole, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		Log(r).Warnf("Unauthorized session delete attempt: %v", err)
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	sessionID, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if *initiatorSession == sessionID {
		Log(r).Debugf("Sessions can't be current")
		httpkit.RenderErr(w, problems.BadRequest(errors.New("session can't be current"))...)
		return
	}

	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	account, err := Domain(r).AccountGet(r.Context(), accountID)
	if err != nil {
		Log(r).Errorf("Failed to get account: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	if identity.CompareRolesUser(*initiatorRole, account.Role) == -1 {
		Log(r).Warn("Account can't delete session of account with higher role")
		httpkit.RenderErr(w, problems.Forbidden("Account can't delete session of account with higher role"))
		return
	}

	err = Domain(r).SessionDelete(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		Log(r).Errorf("Failed to delete device: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	sessions, err := Domain(r).SessionsListByAccount(r.Context(), accountID)
	if err != nil {
		Log(r).Errorf("Failed to retrieve account sessions: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	Log(r).Infof("Sessions Deleted %s for account %s by account %s", sessionID, accountID, initiatorID)
	httpkit.Render(w, responses.SessionCollection(sessions))
}
