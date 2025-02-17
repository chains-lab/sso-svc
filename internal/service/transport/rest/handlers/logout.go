package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/tokens"
)

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	svc := h.svc
	log := svc.Logger

	sessionID, ok := r.Context().Value(tokens.DeviceIDKey).(uuid.UUID)
	if !ok {
		log.Warn("SessionID not found in context")
		httpkit.RenderErr(w, problems.Unauthorized("Sessions not authenticated"))
		return
	}

	userID, ok := r.Context().Value(tokens.UserIDKey).(uuid.UUID)
	if !ok {
		log.Warn("UserID not found in context")
		httpkit.RenderErr(w, problems.Unauthorized("User not authenticated"))
		return
	}

	err := svc.TokenManager.AddToBlackList(r.Context(), userID.String(), sessionID.String())
	if err != nil {
		log.Errorf("Failed to add token to bin: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	err = svc.DB.Sessions.Delete(r, sessionID)
	if err != nil {
		log.Errorf("Failed to delete session: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, http.StatusOK)
}
