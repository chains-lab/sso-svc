package handlers

import (
	"net/http"

	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handlers) AdminUpdateRole(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()

	user, err := tokens.GetUserTokenData(r.Context())
	if err != nil {
		h.presenters.InvalidToken(w, requestID, err)
		return
	}

	updatedUserID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		h.presenters.InvalidParameter(w, requestID, err, "user_id")
		return
	}

	updatedRole, err := roles.ParseRole(chi.URLParam(r, "role"))
	if err != nil {
		h.presenters.InvalidParameter(w, requestID, err, "role")
		return
	}

	appErr := h.app.UpdateUserRole(r.Context(), updatedUserID, updatedRole, user.Role)
	if appErr != nil {
		h.presenters.AppError(w, requestID, appErr)
		return
	}

	h.log.WithField("request_id", requestID).Infof("user %s role updated to %s by %s", updatedUserID, updatedRole, user.UserID)
	httpkit.Render(w, http.StatusOK)
}
