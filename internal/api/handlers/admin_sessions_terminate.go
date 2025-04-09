package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/sso-oauth/internal/app/ape"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) AdminSessionsTerminate(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	if data.AccountID == accountID {
		httpkit.RenderErr(w, problems.Forbidden("You cannot terminate your own session"))
		return
	}

	err = h.app.TerminateByAdmin(r.Context(), accountID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrSessionNotFound):
			h.log.WithError(err).Errorf("error terminating session for account id: %s", accountID)
			httpkit.RenderErr(w, problems.NotFound())
		case errors.Is(err, ape.ErrSessionCannotDeleteForSuperUserByOtherUser):
			h.log.WithError(err).Errorf("error terminating session for account id: %s", accountID)
			httpkit.RenderErr(w, problems.Forbidden("You cannot terminate superuser session"))
		case errors.Is(err, ape.ErrAccountNotFound):
			h.log.WithError(err).Errorf("error terminating session for account id: %s", accountID)
			httpkit.RenderErr(w, problems.NotFound("Account not found"))
		default:
			h.log.WithError(err).Errorf("error terminating session for account id: %s", accountID)
			httpkit.RenderErr(w, problems.InternalError())
		}
		return
	}

	httpkit.Render(w, http.StatusOK)
}
