package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
)

func (h *Handler) OwnAccountGet(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("error getting account data from token")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Error:  err,
		})...)
		return
	}

	res, err := h.app.AccountGetByID(r.Context(), data.AccountID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrAccountDoseNotExits):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusNotFound,
				Title:  "Account not found",
				Detail: "The requested account does not exist.",
			})...)
			return
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}
		h.log.WithError(err).Errorf("error getting account id: %s", data.AccountID)
		return
	}

	httpkit.Render(w, responses.Account(res))
}
