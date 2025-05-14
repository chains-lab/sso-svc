package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) AdminSessionsTerminate(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusUnauthorized,
			Detail: err.Error(),
		})...)
		return
	}

	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:   http.StatusBadRequest,
			Detail:   "invalid account_id format",
			Parametr: "account_id",
		})...)
		return
	}
	if data.AccountID == accountID {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:   http.StatusForbidden,
			Detail:   "You cannot terminate your own session",
			Parametr: "account_id",
		})...)
		return
	}

	err = h.app.TerminateByAdmin(r.Context(), accountID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrSessionDoseNotExits):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Session not found",
				Detail:   "Session does not exist.",
				Parametr: "session_id",
			})...)
		case errors.Is(err, ape.ErrSessionCannotDeleteForSuperUserByOtherUser):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusForbidden,
				Title:    "Forbidden",
				Detail:   "You have not permission to terminate this session.",
				Parametr: "session_id",
			})...)
		case errors.Is(err, ape.ErrAccountDoseNotExits):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Account not found",
				Detail:   "Account does not exist.",
				Parametr: "account_id",
			})...)
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}

		h.log.WithError(err).Errorf("error terminating session for account id: %s", accountID)
		return
	}

	httpkit.Render(w, http.StatusOK)
}
