package handlers

import (
	"net/http"

	"github.com/chains-lab/chains-auth/internal/api/rest/responses"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/google/uuid"
)

func (h *Handler) OwnAccountGet(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()

	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.controllers.TokenData(w, requestID, err)
		return
	}

	res, appErr := h.app.GetAccountByID(r.Context(), user.AccountID)
	if appErr != nil {
		h.controllers.ResultFromApp(w, requestID, appErr)
		return
	}

	httpkit.Render(w, responses.Account(res))
}
