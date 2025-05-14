package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
)

func (h *Handler) SessionsTerminate(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("error getting account data from token")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusUnauthorized,
			Detail: err.Error(),
		})...)
		return
	}

	err = h.app.TerminateByOwner(r.Context(), data.SessionID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrSessionDoseNotExits):
			h.log.WithError(err).Error("session not found session id: %s", data.SessionID)
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusNotFound,
				Title:  "Session not found",
				Detail: "Session does not exist.",
			})...)
			return
		default:
			h.log.WithError(err).Error("error terminating sessions")
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
			return
		}
	}

	httpkit.Render(w, http.StatusNoContent)
}
