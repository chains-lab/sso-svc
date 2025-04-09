package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/sso-oauth/internal/api/responses"
	"github.com/hs-zavet/sso-oauth/internal/app/ape"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) SessionDelete(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	sessionForDeleteId, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiatorSessionID := data.SessionID

	err = h.app.DeleteSessionByOwner(r.Context(), sessionForDeleteId, initiatorSessionID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrSessionNotFound):
			h.log.WithError(err).Errorf("session not found session id: %s", sessionForDeleteId)
			httpkit.RenderErr(w, problems.NotFound())
			return
		default:
			h.log.WithError(err).Errorf("error deleting session")
			httpkit.RenderErr(w, problems.InternalError())
			return
		}
	}

	sessions, err := h.app.GetSessions(r.Context(), data.AccountID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrSessionsNotFound):
			h.log.WithError(err).Error("error getting sessions")
			httpkit.RenderErr(w, problems.InternalError())
			return
		default:
			h.log.WithError(err).Error("error getting sessions")
			httpkit.RenderErr(w, problems.InternalError())
			return
		}
	}

	httpkit.Render(w, responses.SessionCollection(sessions))
}
