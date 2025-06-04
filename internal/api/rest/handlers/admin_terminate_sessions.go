package handlers

import (
	"net/http"

	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handlers) AdminTerminateSessions(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()

	user, err := tokens.GetUserTokenData(r.Context())
	if err != nil {
		h.presenter.InvalidToken(w, requestID, err)
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		h.presenter.InvalidParameter(w, requestID, err, "user_id")
		return
	}

	appErr := h.app.TerminateSessionsByAdmin(r.Context(), userID)
	if appErr != nil {
		h.presenter.AppError(w, requestID, appErr)
		return
	}

	h.log.WithField("request_id", requestID).Infof("Sessions terminated for user %s by admin %s", userID, user.UserID)
	httpkit.Render(w, http.StatusOK)
}
