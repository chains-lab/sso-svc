package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/sso-oauth/internal/service/api/responses"
	"github.com/hs-zavet/sso-oauth/internal/service/domain/ape"
	"github.com/hs-zavet/tokens/identity"
	"github.com/pkg/errors"
)

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		Log(r).Errorf("missing code parameter")
		httpkit.RenderErr(w, problems.BadRequest(errors.New("missing code parameter"))...)
		return
	}

	token, err := GoogleOAuth(r).Exchange(r.Context(), code)
	if err != nil {
		Log(r).WithError(err).Error("failed to exchange code for token")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	client := GoogleOAuth(r).Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			httpkit.RenderErr(w, problems.InternalError())
			return
		}
	}(resp.Body)

	var accountInfo struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&accountInfo); err != nil {
		Log(r).WithError(err).Error("failed to decode account info")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	acc, err := Domain(r).AccountGetByEmail(r.Context(), accountInfo.Email)
	if errors.Is(err, ape.ErrAccountNotFound) {
		tokenAccess, tokenRefresh, err := Domain(r).Login(r.Context(), identity.User, nil, accountInfo.Email, r.UserAgent(), r.RemoteAddr)
		if err != nil {
			Log(r).WithError(err).Error("Failed to login")
			httpkit.RenderErr(w, problems.InternalError())
			return
		}

		httpkit.Render(w, responses.TokensPair(*tokenAccess, *tokenRefresh))
		return
	}

	tokenAccess, tokenRefresh, err := Domain(r).Login(r.Context(), acc.Role, acc.Subscription, acc.Email, r.UserAgent(), r.RemoteAddr)
	if err != nil {
		Log(r).WithError(err).Error("Failed to login")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.TokensPair(*tokenAccess, *tokenRefresh))
}
