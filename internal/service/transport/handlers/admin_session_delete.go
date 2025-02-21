package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/responses"
	"github.com/recovery-flow/tokens"
	"github.com/recovery-flow/tokens/identity"
)

func (h *Handlers) AdminSessionDelete(w http.ResponseWriter, r *http.Request) {
	initiatorID, initiatorSession, initiatorRole, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		h.Log.Warnf("Unauthorized session delete attempt: %v", err)
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	sessionID, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if *initiatorSession == sessionID {
		h.Log.Debugf("Sessions can't be current")
		httpkit.RenderErr(w, problems.BadRequest(errors.New("session can't be current"))...)
		return
	}

	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	account, err := h.Domain.AccountGet(r.Context(), accountID)
	if err != nil {
		h.Log.Errorf("Failed to get account: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	if identity.CompareRolesUser(*initiatorRole, account.Role) == -1 {
		h.Log.Warn("Account can't delete session of account with higher role")
		httpkit.RenderErr(w, problems.Forbidden("Account can't delete session of account with higher role"))
		return
	}

	err = h.Domain.SessionDelete(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		h.Log.Errorf("Failed to delete device: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	sessions, err := h.Domain.SessionsListByAccount(r.Context(), accountID)
	if err != nil {
		h.Log.Errorf("Failed to retrieve account sessions: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	h.Log.Infof("Sessions Deleted %s for account %s by account %s", sessionID, accountID, initiatorID)
	httpkit.Render(w, responses.SessionCollection(sessions))
}
