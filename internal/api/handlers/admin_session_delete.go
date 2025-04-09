package handlers

import (
	"database/sql"
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

func (h *Handler) AdminSessionDelete(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	sessionID, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	err = h.app.DeleteSessionByAdmin(r.Context(), sessionID, data.AccountID, data.SessionID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrAccountNotFound):
			h.log.WithError(err).Errorf("error deleting session %s", sessionID)
			httpkit.RenderErr(w, problems.NotFound("account not found"))
			return
		case errors.Is(err, ape.ErrSessionNotFound):
			h.log.WithError(err).Errorf("error deleting session %s", sessionID)
			httpkit.RenderErr(w, problems.NotFound("session not found"))
			return
		case errors.Is(err, ape.ErrSessionCannotDeleteForSuperUserByOtherUser):
			h.log.WithError(err).Errorf("error deleting session %s", sessionID)
			httpkit.RenderErr(w, problems.Forbidden("session can't be deleted by other user"))
			return
		default:
			h.log.WithError(err).Errorf("error deleting session %s", sessionID)
			httpkit.RenderErr(w, problems.InternalError())
			return
		}
	}

	sessions, err := h.app.GetSessions(r.Context(), accountID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			h.log.WithError(err).Errorf("error deleting session %s", sessionID)
			httpkit.RenderErr(w, problems.NotFound())
			return
		default:
			h.log.WithError(err).Errorf("error getting session for account %s", accountID)
			httpkit.RenderErr(w, problems.InternalError())
			return
		}
	}

	h.log.Infof("delete session %s for account %s by admin: %s", sessionID, accountID, data.AccountID)
	httpkit.Render(w, responses.SessionCollection(sessions))
}
