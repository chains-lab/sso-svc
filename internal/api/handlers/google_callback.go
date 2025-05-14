package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/chains-auth/internal/app"
	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/pkg/errors"
)

func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:   http.StatusBadRequest,
			Detail:   "Code is required.",
			Parametr: "code",
		})...)
		return
	}

	token, err := h.google.Exchange(r.Context(), code)
	if err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusInternalServerError,
		})...)
		h.log.WithError(err).Errorf("error exchanging code for account id: %s", code)
		return
	}

	client := h.google.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusInternalServerError,
		})...)
		h.log.WithError(err).Errorf("error getting account info from google")
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
			h.log.WithError(err).Errorf("error closing response body")
			return
		}
	}(resp.Body)

	var accountInfo struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&accountInfo); err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusInternalServerError,
		})...)
		h.log.WithError(err).Errorf("error decoding account info from google")
		return
	}

	_, err = h.app.AccountGetByEmail(r.Context(), accountInfo.Email)
	if errors.Is(err, ape.ErrAccountDoseNotExits) {
		session, err := h.app.Login(r.Context(), app.LoginRequest{
			Email: accountInfo.Email,
		})
		if err != nil {
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
			h.log.WithError(err).Errorf("error get account for email: %s", accountInfo.Email)
			return
		}

		httpkit.Render(w, responses.TokensPair(session.Access, session.Refresh))
		return
	}

	session, err := h.app.Login(r.Context(), app.LoginRequest{
		Email: accountInfo.Email,
	})
	if err != nil {
		switch {
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}
		h.log.WithError(err).Errorf("error logging in for email: %s", accountInfo.Email)
		return
	}

	httpkit.Render(w, responses.TokensPair(session.Access, session.Refresh))
}
