package handlers

import (
	"net/http"

	"github.com/chains-lab/chains-auth/internal/api/rest/responses"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handlers) AdminDeleteSession(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()

	user, err := tokens.GetUserTokenData(r.Context())
	if err != nil {
		h.presenter.InvalidToken(w, requestID, err)
		return
	}

	sessionID, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		h.presenter.InvalidParameter(w, requestID, err, "session_id")
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		h.presenter.InvalidParameter(w, requestID, err, "user_id")
		return
	}

	appErr := h.app.DeleteSessionByAdmin(r.Context(), sessionID, user.UserID, user.SessionID)
	if appErr != nil {
		h.presenter.AppError(w, requestID, appErr)
		return
	}

	sessions, appErr := h.app.GetUserSessions(r.Context(), userID)
	if appErr != nil {
		h.presenter.AppError(w, requestID, appErr)
		return
	}

	h.log.WithField("request_id", requestID).Infof("delete session %s for user %s by admin: %s", sessionID, userID, user.UserID)
	httpkit.Render(w, responses.SessionCollection(sessions))
}
