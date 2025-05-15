package handlers

import (
	"net/http"

	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) AdminSessionsTerminate(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()

	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.controllers.TokenData(w, requestID, err)
		return
	}

	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		h.controllers.ParameterFromURL(w, requestID, err, "account_id")
		return
	}

	appErr := h.app.TerminateSessionsByAdmin(r.Context(), accountID)
	if appErr != nil {
		h.controllers.ResultFromApp(w, requestID, appErr)
		return
	}

	h.log.WithField("request_id", requestID).Infof("Sessions terminated for account %s by admin %s", accountID, user.AccountID)
	httpkit.Render(w, http.StatusOK)
}
