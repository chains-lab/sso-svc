package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/rest/responses"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
)

func (h *Handlers) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()
	log := h.log.WithField("request_id", requestID)

	code := r.URL.Query().Get("code")
	if code == "" {
		h.presenter.InvalidQuery(w, requestID, "code", fmt.Errorf("missing code"))
		return
	}

	token, err := h.google.Exchange(r.Context(), code)
	if err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusInternalServerError,
		})...)
		log.WithError(err).Errorf("error exchanging code for user id: %s", code)
		return
	}

	client := h.google.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusInternalServerError,
		})...)
		log.WithError(err).Errorf("error getting user info from google")
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

	var userInfo struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusInternalServerError,
		})...)
		log.WithError(err).Errorf("error decoding user info from google")
		return
	}

	_, tokensPair, appErr := h.app.Login(r.Context(), userInfo.Email, roles.User, r.Header.Get("User-Agent"))
	if appErr != nil {
		h.presenter.AppError(w, requestID, appErr)
	}

	h.log.Infof("User %s logged in with Google", userInfo.Email)
	httpkit.Render(w, responses.TokensPair(tokensPair.Access, tokensPair.Refresh))
}
