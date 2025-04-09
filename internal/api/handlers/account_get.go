package handlers

import (
	"errors"
	"net/http"

	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/sso-oauth/internal/api/responses"
	"github.com/hs-zavet/sso-oauth/internal/app/ape"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) AccountGet(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("error getting account data from token")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	res, err := h.app.AccountGetByID(r.Context(), data.AccountID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrAccountNotFound):
			h.log.WithError(err).Errorf("account id: %s", data.AccountID)
			httpkit.RenderErr(w, problems.NotFound("account not found"))
			return
		default:
			h.log.WithError(err).Errorf("error getting account id: %s", data.AccountID)
			httpkit.RenderErr(w, problems.InternalError())
			return
		}
	}

	httpkit.Render(w, responses.Account(res))
}
