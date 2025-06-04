package handlers

import (
	"net/http"

	"github.com/chains-lab/chains-auth/internal/api/rest/responses"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handlers) AdminGetSessions(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		h.presenter.InvalidParameter(w, requestID, err, "user_id")
		return
	}

	sessions, appErr := h.app.GetUserSessions(r.Context(), userID)
	if appErr != nil {
		h.presenter.AppError(w, requestID, appErr)
		return
	}

	httpkit.Render(w, responses.SessionCollection(sessions))
}
