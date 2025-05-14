package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) SessionDelete(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("error getting account data from token")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusUnauthorized,
			Detail: err.Error(),
		})...)
		return
	}

	sessionForDeleteId, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:   http.StatusBadRequest,
			Detail:   "Session ID must be a valid UUID.",
			Parametr: "session_id",
		})...)
		return
	}

	initiatorSessionID := data.SessionID

	err = h.app.DeleteSessionByOwner(r.Context(), sessionForDeleteId, initiatorSessionID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrSessionDoseNotExits):
			h.log.WithError(err).Errorf("session not found session id: %s", sessionForDeleteId)
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusNotFound,
				Title:  "Session not found",
				Detail: "The requested session does not exist.",
			})...)
			return
		default:
			h.log.WithError(err).Errorf("error deleting session")
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
			return
		}
	}

	sessions, err := h.app.GetSessions(r.Context(), data.AccountID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrSessionsForAccountDoseNotExits):
			h.log.WithError(err).Error("error getting sessions")
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusNotFound,
				Title:  "Sessions not found",
				Detail: "The requested sessions do not exist.",
			})...)
			return
		default:
			h.log.WithError(err).Error("error getting sessions")
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
			return
		}
	}

	httpkit.Render(w, responses.SessionCollection(sessions))
}
