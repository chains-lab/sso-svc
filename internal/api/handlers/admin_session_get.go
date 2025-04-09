package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/sso-oauth/internal/api/responses"
	"github.com/hs-zavet/sso-oauth/internal/app/ape"
)

func (h *Handler) AdminSessionGet(w http.ResponseWriter, r *http.Request) {
	sessionID, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"session_id": validation.NewError("session_id", "Invalid session_id"),
		})...)
		return
	}

	session, err := h.app.GetSession(r.Context(), sessionID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrSessionNotFound):
			h.log.WithError(err).Errorf("session id: %s", sessionID)
			httpkit.RenderErr(w, problems.NotFound())
			return
		default:
			h.log.WithError(err).Errorf("error getting session id: %s", sessionID)
			httpkit.RenderErr(w, problems.InternalError())
			return
		}
	}

	httpkit.Render(w, responses.Session(session))
}
