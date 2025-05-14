package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) AdminSessionDelete(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("error getting account data from token")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusUnauthorized,
			Detail: err.Error(),
		})...)
		return
	}

	sessionID, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:   http.StatusBadRequest,
			Detail:   "Session ID must be a valid UUID.",
			Parametr: "session_id",
		})...)
		return
	}

	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:   http.StatusBadRequest,
			Detail:   "Account ID must be a valid UUID.",
			Parametr: "account_id",
		})...)
		return
	}

	err = h.app.DeleteSessionByAdmin(r.Context(), sessionID, data.AccountID, data.SessionID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrAccountDoseNotExits):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Account not found",
				Detail:   "The requested account does not exist.",
				Parametr: "account_id",
			})...)
		case errors.Is(err, ape.ErrSessionDoseNotExits):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Session not found",
				Detail:   "The requested session does not exist.",
				Parametr: "session_id",
			})...)
		case errors.Is(err, ape.ErrSessionCannotDeleteForSuperUserByOtherUser):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusForbidden,
				Title:  "Forbidden",
				Detail: "You have not permission to delete this session.",
			})...)
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}

		h.log.WithError(err).Errorf("error deleting session %s", sessionID)
		return
	}

	sessions, err := h.app.GetSessions(r.Context(), accountID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			h.log.WithError(err).Errorf("error deleting session %s", sessionID)
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusNotFound,
				Title:  "Session not found",
			})...)
			return
		default:
			h.log.WithError(err).Errorf("error getting session for account %s", accountID)
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
			return
		}
	}

	h.log.Infof("delete session %s for account %s by admin: %s", sessionID, accountID, data.AccountID)
	httpkit.Render(w, responses.SessionCollection(sessions))
}
