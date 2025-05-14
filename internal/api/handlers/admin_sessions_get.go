package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) AdminSessionsGet(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:   http.StatusBadRequest,
			Detail:   "Account ID must be a valid UUID.",
			Parametr: "account_id",
		})...)
		return
	}

	sessions, err := h.app.GetSessions(r.Context(), accountID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrSessionsForAccountDoseNotExits):
			h.log.WithError(err).Errorf("account id: %s", accountID)
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Sessions not found",
				Detail:   "Session for this account does not exist.",
				Parametr: "session_id",
			})...)
			return
		case errors.Is(err, ape.ErrAccountDoseNotExits):
			h.log.WithError(err).Errorf("account id: %s", accountID)
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Account not found",
				Detail:   "Account does not exist.",
				Parametr: "account_id",
			})...)
			return
		default:
			h.log.WithError(err).Errorf("error getting sessions for account %s", accountID)
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
			return
		}
	}

	httpkit.Render(w, responses.SessionCollection(sessions))
}
