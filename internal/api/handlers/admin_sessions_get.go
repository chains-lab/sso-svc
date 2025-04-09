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

func (h *Handler) AdminSessionsGet(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"account_id": validation.NewError("account_id", "Invalid account_id"),
		})...)
		return
	}

	sessions, err := h.app.GetSessions(r.Context(), accountID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrSessionNotFound):
			h.log.WithError(err).Errorf("account id: %s", accountID)
			httpkit.RenderErr(w, problems.NotFound())
			return
		case errors.Is(err, ape.ErrAccountNotFound):
			h.log.WithError(err).Errorf("account id: %s", accountID)
			httpkit.RenderErr(w, problems.NotFound("account not found"))
			return
		default:
			h.log.WithError(err).Errorf("error getting sessions for account %s", accountID)
			httpkit.RenderErr(w, problems.InternalError())
			return
		}
	}

	httpkit.Render(w, responses.SessionCollection(sessions))
}
