package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/tokens"
	"github.com/hs-zavet/tokens/identity"
)

func (h *Handler) AdminSessionsTerminate(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountData(r.Context())
	if err != nil {
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	account, err := h.app.AccountGetByID(r.Context(), accountID)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}
	//
	//accRole, err := identity.ParseIdentityType(account.Role)
	//if err != nil {
	//	httpkit.RenderErr(w, problems.InternalError())
	//	return
	//}

	if identity.CompareRolesUser(data.Role, account.Role) == -1 {
		httpkit.RenderErr(w, problems.Forbidden("Account can't delete session of account with higher role"))
		return
	}

	err = h.app.Terminate(r.Context(), accountID)
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
