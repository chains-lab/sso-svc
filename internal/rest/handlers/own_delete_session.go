package handlers

import (
	"net/http"

	"github.com/chains-lab/chains-auth/internal/rest/responses"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handlers) DeleteSession(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()

	user, err := tokens.GetUserTokenData(r.Context())
	if err != nil {
		h.presenter.InvalidToken(w, requestID, err)
		return
	}

	sessionForDeleteId, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		h.presenter.InvalidParameter(w, requestID, err, "session_id")
		return
	}

	initiatorSessionID := user.SessionID

	appErr := h.app.DeleteSessionByOwner(r.Context(), sessionForDeleteId, initiatorSessionID)
	if appErr != nil {
		h.presenter.AppError(w, requestID, appErr)
		return
	}

	sessions, appErr := h.app.GetUserSessions(r.Context(), user.UserID)
	if appErr != nil {
		h.presenter.AppError(w, requestID, appErr)
		return
	}

	h.log.WithField("request_id", requestID).Infof("logout session id: %s", user.SessionID)
	httpkit.Render(w, responses.SessionCollection(sessions))
}
