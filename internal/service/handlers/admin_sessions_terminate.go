package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/cifractx"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/roles"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/tokens"
	"github.com/sirupsen/logrus"
)

func AdminSessionsTerminate(w http.ResponseWriter, r *http.Request) {
	server, err := cifractx.GetValue[*config.Server](r.Context(), config.SERVER)
	if err != nil {
		logrus.Errorf("Failed to retrieve service configuration %s", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}
	log := server.Logger

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

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	user, err := server.SqlDB.Accounts.GetById(r, userID)
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
		log.Warn("User can't terminate sessions of higher level user")
		httpkit.RenderErr(w, problems.Forbidden("User can't terminate sessions of higher level user"))
		return
	}

	err = server.SqlDB.Sessions.TerminateSessions(r, userID, nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		httpkit.RenderErr(w, problems.InternalError("Failed to terminate sessions"))
		return
	}

	log.Infof("Session terminated for user %s by user %s", userID, initiatorID)
	httpkit.Render(w, http.StatusOK)
}
