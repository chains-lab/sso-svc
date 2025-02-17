package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/tokens"
)

func (h *Handlers) SessionsTerminate(w http.ResponseWriter, r *http.Request) {
	svc := h.svc
	log := svc.Logger

	userID, ok := r.Context().Value(tokens.UserIDKey).(uuid.UUID)
	if !ok {
		log.Warn("UserID not found in context")
		httpkit.RenderErr(w, problems.Unauthorized("User not authenticated"))
		return
	}

	sessionID, ok := r.Context().Value(tokens.DeviceIDKey).(uuid.UUID)
	if !ok {
		log.Warn("DeviceID not found in context")
		httpkit.RenderErr(w, problems.Unauthorized("Device not authenticated"))
		return
	}

	sessions, err := svc.DB.Sessions.SelectByUserID(r, userID)

	for _, ses := range sessions {
		err = svc.TokenManager.AddToBlackList(r.Context(), userID.String(), ses.ID.String())
		if err != nil {
			log.Errorf("Failed to add token to bin: %v", err)
			httpkit.RenderErr(w, problems.InternalError())
			return
		}
	}

	err = svc.DB.Sessions.Terminate(r, userID, &sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		httpkit.RenderErr(w, problems.InternalError("Failed to terminate sessions"))
		return
	}

	httpkit.Render(w, http.StatusOK)
}
