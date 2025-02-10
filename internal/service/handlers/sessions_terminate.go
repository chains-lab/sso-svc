package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/cifractx"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/tokens"
)

func SessionsTerminate(w http.ResponseWriter, r *http.Request) {
	Server, err := cifractx.GetValue[*config.Server](r.Context(), config.SERVER)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError("Failed to retrieve service configuration"))
		return
	}

	log := Server.Logger

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

	sessions, err := Server.SqlDB.Sessions.GetSessions(r, userID)

	for _, ses := range sessions {
		err = Server.TokenManager.Bin.Add(userID.String(), ses.ID.String())
		if err != nil {
			log.Errorf("Failed to add token to bin: %v", err)
			httpkit.RenderErr(w, problems.InternalError())
			return
		}
	}

	err = Server.SqlDB.Sessions.TerminateSessions(r, userID, &sessionID)
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
