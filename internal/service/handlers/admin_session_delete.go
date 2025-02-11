package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/roles"
	"github.com/recovery-flow/sso-oauth/internal/service/responses"
	"github.com/recovery-flow/tokens"
)

func (h *Handlers) AdminSessionDelete(w http.ResponseWriter, r *http.Request) {
	svc := h.svc
	log := svc.Logger

	initiatorID, ok := r.Context().Value(tokens.UserIDKey).(uuid.UUID)
	if !ok {
		log.Warn("UserID not found in context")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	InitiatorRoleStr, ok := r.Context().Value(tokens.RoleKey).(string)
	if !ok {
		log.Warn("Role not found in context")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	InitiatorRole, err := roles.StringToRoleUser(InitiatorRoleStr)
	if err != nil {
		log.Errorf("Failed to parse Initiator updatedRole: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	initiatorSessionID, ok := r.Context().Value(tokens.DeviceIDKey).(uuid.UUID)
	if !ok {
		log.Warn("SessionID not found in context")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	sessionID, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if initiatorSessionID == sessionID {
		log.Debugf("Session can't be current")
		httpkit.RenderErr(w, problems.BadRequest(errors.New("session can't be current"))...)
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	user, err := svc.SqlDB.Accounts.GetById(r, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		log.Errorf("Failed to retrieve user: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	userRole, err := roles.StringToRoleUser(user.Role)
	if err != nil {
		log.Errorf("Failed to parse user role: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	if roles.CompareRolesUser(InitiatorRole, userRole) == -1 {
		log.Warn("User can't delete session of user with higher role")
		httpkit.RenderErr(w, problems.Forbidden("User can't delete session of user with higher role"))
		return
	}

	err = svc.SqlDB.Sessions.Delete(r, sessionID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		log.Errorf("Failed to delete device: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	err = svc.TokenManager.Bin.Add(userID.String(), sessionID.String())
	if err != nil {
		log.Errorf("Failed to add token to bin: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	sessions, err := svc.SqlDB.Sessions.GetSessions(r, userID)
	if err != nil {
		log.Errorf("Failed to retrieve user sessions: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	log.Infof("Session Dleted %s for user %s by user %s", sessionID, userID, initiatorID)
	httpkit.Render(w, responses.SessionCollection(sessions))
}
