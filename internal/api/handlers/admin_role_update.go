package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/tokens"
	"github.com/hs-zavet/tokens/identity"
)

func (h *Handler) AdminRoleUpdate(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountData(r.Context())
	if err != nil {
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	updatedAccountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"account_id": validation.NewError("account_id", "invalid UUID"),
		})...)
		return
	}

	updatedRole, err := identity.ParseIdentityType(chi.URLParam(r, "role"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"role": validation.NewError("role", "invalid role"),
		})...)
		return
	}

	if identity.CompareRolesUser(data.Role, updatedRole) != 1 {
		httpkit.RenderErr(w, problems.Forbidden("Account can't update role to higher level"))
		return
	}

	err = h.app.AccountUpdateRole(r.Context(), updatedAccountID, string(updatedRole))
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, http.StatusOK)
}
