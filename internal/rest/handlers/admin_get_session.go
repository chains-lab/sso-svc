package handlers

import (
	"net/http"

	"github.com/chains-lab/chains-auth/internal/rest/responses"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handlers) AdminGetSession(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()

	sessionID, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		h.presenters.InvalidParameter(w, requestID, err, "session_id")
		return
	}

	session, appErr := h.app.GetSession(r.Context(), sessionID)
	if appErr != nil {
		h.presenters.AppError(w, requestID, appErr)
		return
	}

	httpkit.Render(w, responses.Session(session))
}
