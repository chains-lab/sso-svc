package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/api/responses"
	"github.com/recovery-flow/tokens/identity"
)

func (h *Handlers) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		h.Log.Errorf("missing code parameter")
		httpkit.RenderErr(w, problems.BadRequest(errors.New("missing code parameter"))...)
		return
	}

	token, err := h.GoogleOAuth.Exchange(r.Context(), code)
	if err != nil {
		h.Log.WithError(err).Error("failed to exchange code for token")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	client := h.GoogleOAuth.Client(r.Context(), token)
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
		h.Log.WithError(err).Error("failed to decode account info")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	tokenAccess, tokenRefresh, err := h.Domain.Login(r.Context(), identity.User, accountInfo.Email, r.UserAgent(), r.RemoteAddr)
	if err != nil {
		h.Log.WithError(err).Error("Failed to login")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.TokensPair(*tokenAccess, *tokenRefresh))
}
