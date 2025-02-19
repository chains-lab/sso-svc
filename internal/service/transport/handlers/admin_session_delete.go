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
	"github.com/recovery-flow/sso-oauth/internal/service/transport/responses"
	"github.com/recovery-flow/tokens"
)

func (a *App) AdminSessionDelete(w http.ResponseWriter, r *http.Request) {
	initiatorID, initiatorSession, initiatorRole, err := tokens.GetAccountData(r.Context())
	if err != nil {
		a.Log.Warnf("Unauthorized session delete attempt: %v", err)
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	sessionID, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if *initiatorSession == sessionID {
		a.Log.Debugf("Sessions can't be current")
		httpkit.RenderErr(w, problems.BadRequest(errors.New("session can't be current"))...)
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	user, err := a.Domain.AccountGet(r.Context(), userID)
	if err != nil {
		a.Log.Errorf("Failed to get user: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	if roles.CompareRolesUser(*initiatorRole, user.Role) == -1 {
		a.Log.Warn("User can't delete session of user with higher role")
		httpkit.RenderErr(w, problems.Forbidden("User can't delete session of user with higher role"))
		return
	}

	err = a.Domain.SessionDelete(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		a.Log.Errorf("Failed to delete device: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	sessions, err := a.Domain.SessionsListByUser(r.Context(), userID)
	if err != nil {
		a.Log.Errorf("Failed to retrieve user sessions: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	a.Log.Infof("Sessions Deleted %s for user %s by user %s", sessionID, userID, initiatorID)
	httpkit.Render(w, responses.SessionCollection(sessions))
}
