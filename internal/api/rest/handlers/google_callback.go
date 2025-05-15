package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/api/rest/responses"
	"github.com/chains-lab/chains-auth/internal/app"
	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/google/uuid"
)

func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()
	log := h.log.WithField("request_id", requestID)

	code := r.URL.Query().Get("code")
	if code == "" {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:    http.StatusBadRequest,
			Code:      ape.CodeInvalidRequestQuery,
			Title:     "Invalid request query",
			Detail:    "Code is required.",
			Parameter: "code",
		})...)
		return
	}

	token, err := h.google.Exchange(r.Context(), code)
	if err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusInternalServerError,
		})...)
		log.WithError(err).Errorf("error exchanging code for account id: %s", code)
		return
	}

	client := h.google.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusInternalServerError,
		})...)
		log.WithError(err).Errorf("error getting account info from google")
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
			log.WithError(err).Errorf("error closing response body")
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
		log.WithError(err).Errorf("error decoding account info from google")
		return
	}

	_, appErr := h.app.GetAccountByEmail(r.Context(), accountInfo.Email)
	if appErr.Code == ape.CodeAccountDoesNotExist {
		session, appErr := h.app.Login(r.Context(), app.LoginRequest{
			Email: accountInfo.Email,
		})
		if appErr != nil {
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
			log.WithError(appErr.Unwrap()).Errorf("error get account for email: %s", accountInfo.Email)
			return
		}

		httpkit.Render(w, responses.TokensPair(session.Access, session.Refresh))
		return
	}

	session, appErr := h.app.Login(r.Context(), app.LoginRequest{
		Email: accountInfo.Email,
	})
	if appErr != nil {
		switch {
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}
		log.WithError(appErr.Unwrap()).Errorf("error logging in for email: %s", accountInfo.Email)
		return
	}

	httpkit.Render(w, responses.TokensPair(session.Access, session.Refresh))
}
