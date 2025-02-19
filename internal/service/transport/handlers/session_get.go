package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/responses"
	"github.com/recovery-flow/tokens"
)

func (a *App) SessionGet(w http.ResponseWriter, r *http.Request) {
	accountID, _, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		a.Log.Warnf("Unauthorized session get attempt: %v", err)
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	sessionID, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		a.Log.Errorf("Failed to parse session_id: %v", err)
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"session_id": validation.NewError("session_id", "Invalid session_id"),
		})...)
		return
	}

	session, err := a.Domain.SessionGet(r.Context(), sessionID)
	if err != nil {
		a.Log.Errorf("Failed to retrieve user session: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	if session.UserID != *accountID {
		a.Log.Debugf("Session doesn't belong to user")
		httpkit.RenderErr(w, problems.Forbidden("Session doesn't belong to user"))
		return
	}

	httpkit.Render(w, responses.Session(*session))
}
