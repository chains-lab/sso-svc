package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/sso-oauth/internal/api/responses"
	"github.com/hs-zavet/sso-oauth/internal/app/ape"
)

func (h *Handler) AdminAccountGet(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid account_id"))...)
		return
	}

	res, err := h.app.AccountGetByID(r.Context(), accountID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrAccountNotFound):
			h.log.WithError(err).Errorf("account id: %s", accountID)
			httpkit.RenderErr(w, problems.NotFound("account not found"))
			return
		default:
			h.log.WithError(err).Errorf("error getting account id: %s", accountID)
			httpkit.RenderErr(w, problems.InternalError())
			return
		}
	}
	httpkit.Render(w, responses.Account(res))
}
