package handlers

import (
	"net/http"

	"github.com/chains-lab/chains-auth/internal/api/rest/responses"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) AdminGetSessions(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()

	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		h.controllers.ParameterFromURL(w, requestID, err, "account_id")
		return
	}

	sessions, appErr := h.app.GetSessions(r.Context(), accountID)
	if appErr != nil {
		h.controllers.ResultFromApp(w, requestID, appErr)
		return
	}

	httpkit.Render(w, responses.SessionCollection(sessions))
}
