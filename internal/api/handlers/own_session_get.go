package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
)

func (h *Handler) OwnSessionGet(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("error getting account data from token")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusUnauthorized,
			Detail: err.Error(),
		})...)
		return
	}

	session, err := h.app.GetSession(r.Context(), data.SessionID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrSessionDoseNotExits):
			h.log.WithError(err).Errorf("session not found session id: %s", data.SessionID)
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusNotFound,
				Title:  "Session not found",
				Detail: "Session does not exist.",
			})...)
			return
		default:
			h.log.WithError(err).Errorf("error getting session")
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
			return
		}
	}

	httpkit.Render(w, responses.Session(session))
}
