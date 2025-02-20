package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/responses"
	"github.com/recovery-flow/tokens/identity"
)

func (a *App) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		httpkit.RenderErr(w, problems.BadRequest(errors.New("missing code parameter"))...)
		return
	}

	token, err := a.GoogleOAuth.Exchange(r.Context(), code)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	client := a.GoogleOAuth.Client(r.Context(), token)
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

	var userInfo struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		a.Log.Errorf("failed to decode user info: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	tokenAccess, tokenRefresh, err := a.Domain.SessionLogin(r.Context(), identity.User, userInfo.Email, r.UserAgent(), r.RemoteAddr)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.TokensPair(tokenAccess, tokenRefresh))
}
