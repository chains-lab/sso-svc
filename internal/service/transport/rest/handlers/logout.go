package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/tokens"
)

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {

	sessionID, ok := r.Context().Value(tokens.DeviceIDKey).(uuid.UUID)
	if !ok {
		httpkit.RenderErr(w, problems.Unauthorized("Sessions not authenticated"))
		return
	}

	userID, ok := r.Context().Value(tokens.UserIDKey).(uuid.UUID)
	if !ok {
		httpkit.RenderErr(w, problems.Unauthorized("User not authenticated"))
		return
	}

	err := h.svc.TokenManager.AddToBlackList(r.Context(), sessionID.String(), userID.String())
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	err = h.svc.Domain.Session.Logout(r.Context(), sessionID.String(), userID.String())
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, http.StatusOK)
}
