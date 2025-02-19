package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/roles"
	"github.com/recovery-flow/tokens"
)

func (a *App) AdminRoleUpdate(w http.ResponseWriter, r *http.Request) {
	initiatorID, _, InitiatorRole, err := tokens.GetAccountData(r.Context())
	if err != nil {
		a.Log.Warnf("Unauthorized role update attempt: %v", err)
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	updatedUserID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"user_id": validation.NewError("user_id", "invalid UUID"),
		})...)
		return
	}

	updatedRole, err := roles.ParseUserRole(chi.URLParam(r, "role"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"role": validation.NewError("role", "invalid role"),
		})...)
	}

	if roles.CompareRolesUser(*InitiatorRole, updatedRole) != 1 {
		a.Log.Warn("User can't update role to higher level than his own")
		httpkit.RenderErr(w, problems.Forbidden("User can't update role to higher level"))
		return
	}

	user, err := a.Domain.AccountGet(r.Context(), updatedUserID)
	if err != nil {
		a.Log.Errorf("Failed to get user: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	if roles.CompareRolesUser(*InitiatorRole, user.Role) == -1 {
		a.Log.Warn("User can't update role of user with higher role than his own")
		httpkit.RenderErr(w, problems.Forbidden("User can't update role of user with higher role"))
		return
	}

	_, err = a.Domain.AccountUpdateRole(r.Context(), updatedUserID, string(updatedRole))
	if err != nil {
		a.Log.Errorf("Failed to update role: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	a.Log.Infof("Role updated for user %s to %s by user %s", updatedUserID, updatedRole, initiatorID)
	httpkit.Render(w, http.StatusOK)
}
