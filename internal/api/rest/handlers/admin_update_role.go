package handlers

import (
	"net/http"

	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) AdminRoleUpdate(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()

	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.controllers.TokenData(w, requestID, err)
		return
	}

	updatedAccountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		h.controllers.ParameterFromURL(w, requestID, err, "account_id")
		return
	}

	updatedRole, err := roles.ParseRole(chi.URLParam(r, "role"))
	if err != nil {
		h.controllers.ParameterFromURL(w, requestID, err, "role")
		return
	}

	appErr := h.app.UpdateAccountRole(r.Context(), updatedAccountID, updatedRole, user.Role)
	if appErr != nil {
		h.controllers.ResultFromApp(w, requestID, appErr)
		return
	}

	h.log.WithField("request_id", requestID).Infof("account %s role updated to %s by %s", updatedAccountID, updatedRole, user.AccountID)
	httpkit.Render(w, http.StatusOK)
}
