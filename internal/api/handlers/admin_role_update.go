package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) AdminRoleUpdate(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("error getting account data from token")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusUnauthorized,
			Detail: err.Error(),
		})...)
		return
	}

	updatedAccountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		h.log.WithError(err).Error("error parsing account_id")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:   http.StatusBadRequest,
			Detail:   "Account ID must be a valid UUID.",
			Parametr: "account_id",
		})...)
		return
	}

	updatedRole, err := roles.ParseRole(chi.URLParam(r, "role"))
	if err != nil {
		h.log.WithError(err).Error("error parsing role")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:   http.StatusBadRequest,
			Detail:   "Role must be a valid role.",
			Parametr: "role",
		})...)
		return
	}

	err = h.app.AccountUpdateRole(r.Context(), updatedAccountID, updatedRole, data.Role)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrAccountDoseNotExits):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Account not found",
				Detail:   "Account does not exist.",
				Parametr: "account_id",
			})...)
		case errors.Is(err, ape.ErrUserHasNoPermissionToUpdateRole):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusForbidden,
				Detail: "You do not have permission to update this account's role.",
			})...)
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}
	}

	h.log.Infof("account %s role updated to %s by %s", updatedAccountID, updatedRole, data.AccountID)
	httpkit.Render(w, http.StatusOK)
}
