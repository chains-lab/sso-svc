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

func (h *Handlers) AdminSessionsTerminate(w http.ResponseWriter, r *http.Request) {
	initiatorID, _, initiatorRole, _, err := tokens.GetAccountData(r.Context())
	userID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		h.Log.WithError(err).Warn("Invalid account_id")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	account, err := h.Domain.AccountGet(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	if identity.CompareRolesUser(*initiatorRole, account.Role) == -1 {
		h.Log.Errorf("User can't terminate sessions of higher level account")
		httpkit.RenderErr(w, problems.Forbidden("User can't terminate sessions of higher level account"))
		return
	}

	err = h.Domain.SessionsTerminate(r.Context(), userID, nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		httpkit.RenderErr(w, problems.InternalError("Failed to terminate sessions"))
		return
	}

	h.Log.Infof("Sessions terminated for account %s by account %s", userID, initiatorID)
	httpkit.Render(w, http.StatusOK)
}
