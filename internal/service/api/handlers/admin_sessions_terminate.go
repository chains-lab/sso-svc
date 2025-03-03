package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/tokens"
	"github.com/recovery-flow/tokens/identity"
)

func AdminSessionsTerminate(w http.ResponseWriter, r *http.Request) {
	initiatorID, _, _, initiatorRole, _, err := tokens.GetAccountData(r.Context())
	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		Log(r).WithError(err).Warn("Invalid account_id")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	account, err := Domain(r).AccountGet(r.Context(), accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	if identity.CompareRolesUser(*initiatorRole, account.Role) == -1 {
		Log(r).Errorf("Account can't terminate sessions of higher level account")
		httpkit.RenderErr(w, problems.Forbidden("Account can't terminate sessions of higher level account"))
		return
	}

	err = Domain(r).SessionsTerminate(r.Context(), accountID, nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		httpkit.RenderErr(w, problems.InternalError("Failed to terminate sessions"))
		return
	}

	Log(r).Infof("Sessions terminated for account %s by account %s", accountID, initiatorID)
	httpkit.Render(w, http.StatusOK)
}
