package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/tokens"
	"github.com/recovery-flow/tokens/identity"
)

func (a *App) AdminSessionsTerminate(w http.ResponseWriter, r *http.Request) {
	initiatorID, _, initiatorRole, _, err := tokens.GetAccountData(r.Context())
	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	account, err := a.Domain.AccountGet(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		a.Log.Errorf("Failed to retrieve account: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	if identity.CompareRolesUser(*initiatorRole, account.Role) == -1 {
		a.Log.Warn("User can't terminate sessions of higher level account")
		httpkit.RenderErr(w, problems.Forbidden("User can't terminate sessions of higher level account"))
		return
	}

	err = a.Domain.SessionsTerminate(r.Context(), userID, nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		httpkit.RenderErr(w, problems.InternalError("Failed to terminate sessions"))
		return
	}

	a.Log.Infof("Sessions terminated for account %s by account %s", userID, initiatorID)
	httpkit.Render(w, http.StatusOK)
}
