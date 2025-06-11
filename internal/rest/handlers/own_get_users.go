package handlers

import (
	"net/http"

	"github.com/chains-lab/chains-auth/internal/rest/responses"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/google/uuid"
)

func (h *Handlers) OwnUserGet(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()

	user, err := tokens.GetUserTokenData(r.Context())
	if err != nil {
		h.presenters.InvalidToken(w, requestID, err)
		return
	}

	res, appErr := h.app.GetUserByID(r.Context(), user.UserID)
	if appErr != nil {
		h.presenters.AppError(w, requestID, appErr)
		return
	}

	httpkit.Render(w, responses.User(res))
}
