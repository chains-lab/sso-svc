package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/tokens"
	"github.com/recovery-flow/tokens/identity"
)

func (h *Handlers) AdminRoleUpdate(w http.ResponseWriter, r *http.Request) {
	initiatorID, _, InitiatorRole, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		h.Log.WithError(err).Warn("Unauthorized role update attempt")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	updatedAccountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		h.Log.WithError(err).Warn("Invalid account_id")
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"account_id": validation.NewError("account_id", "invalid UUID"),
		})...)
		return
	}

	updatedRole, err := identity.ParseIdentityType(chi.URLParam(r, "role"))
	if err != nil {
		h.Log.WithError(err).Warn("Invalid role")
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"role": validation.NewError("role", "invalid role"),
		})...)
		return
	}

	if identity.CompareRolesUser(*InitiatorRole, updatedRole) != 1 {
		h.Log.Warn("Account can't update role to higher level than his own")
		httpkit.RenderErr(w, problems.Forbidden("Account can't update role to higher level"))
		return
	}

	account, err := h.Domain.AccountGet(r.Context(), updatedAccountID)
	if err != nil {
		h.Log.WithError(err).Warn("Failed to get account")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	if identity.CompareRolesUser(*InitiatorRole, account.Role) == -1 {
		h.Log.Warn("Account can't update role of account with higher role than his own")
		httpkit.RenderErr(w, problems.Forbidden("Account can't update role of account with higher role"))
		return
	}

	_, err = h.Domain.AccountUpdateRole(r.Context(), updatedAccountID, updatedRole)
	if err != nil {
		h.Log.WithError(err).Warn("Failed to update role")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	h.Log.Infof("Role updated for account %s to %s by account %s", updatedAccountID, updatedRole, initiatorID)
	httpkit.Render(w, http.StatusOK)
}
