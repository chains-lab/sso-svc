package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/sso-oauth/internal/app/ape"
	"github.com/hs-zavet/tokens"
	"github.com/hs-zavet/tokens/roles"
)

func (h *Handler) AdminRoleUpdate(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	updatedAccountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		h.log.WithError(err).Error("error parsing account_id")
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"account_id": validation.NewError("account_id", "invalid UUID"),
		})...)
		return
	}

	updatedRole, err := roles.ParseRole(chi.URLParam(r, "role"))
	if err != nil {
		h.log.WithError(err).Error("error parsing role")
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"role": validation.NewError("role", "invalid role"),
		})...)
		return
	}

	err = h.app.AccountUpdateRole(r.Context(), updatedAccountID, updatedRole, data.Role)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrAccountNotFound):
			h.log.WithError(err).Errorf("account id: %s", updatedAccountID)
			httpkit.RenderErr(w, problems.NotFound("account not found"))
			return
		default:
			h.log.WithError(err).Errorf("error updating role for account id: %s", updatedAccountID)
			httpkit.RenderErr(w, problems.InternalError())
			return
		}
	}

	h.log.Infof("account %s role updated to %s by %s", updatedAccountID, updatedRole, data.AccountID)
	httpkit.Render(w, http.StatusOK)
}
