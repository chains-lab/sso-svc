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

func (h *Handler) SessionGet(w http.ResponseWriter, r *http.Request) {
	svc := h.svc
	log := svc.Logger

	userID, ok := r.Context().Value(tokens.UserIDKey).(uuid.UUID)
	if !ok {
		log.Warn("UserID not found in context")
		httpkit.RenderErr(w, problems.Unauthorized("User not authenticated"))
		return
	}

	sessionID, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		log.Errorf("Failed to parse session_id: %v", err)
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"session_id": validation.NewError("session_id", "Invalid session_id"),
		})...)
		return
	}

	session, err := svc.DB.Sessions.GetByID(r, sessionID)
	if err != nil {
		log.Errorf("Failed to retrieve user session: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	if session.UserID != userID {
		log.Debugf("Session doesn't belong to user")
		httpkit.RenderErr(w, problems.Forbidden("Session doesn't belong to user"))
		return
	}

	httpkit.Render(w, responses.Session(*session))
}
