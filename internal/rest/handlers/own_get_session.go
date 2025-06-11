package handlers

import (
	"net/http"

	"github.com/chains-lab/chains-auth/internal/rest/responses"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/google/uuid"
)

func (h *Handlers) OwnGetSession(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()

	user, err := tokens.GetUserTokenData(r.Context())
	if err != nil {
		h.presenters.InvalidToken(w, requestID, err)
		return
	}

	session, appErr := h.app.GetSession(r.Context(), user.SessionID)
	if appErr != nil {
		h.presenters.AppError(w, requestID, appErr)
		return
	}

	httpkit.Render(w, responses.Session(session))
}
