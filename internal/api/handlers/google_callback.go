package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/sso-oauth/internal/api/responses"
	"github.com/hs-zavet/sso-oauth/internal/app"
	"github.com/hs-zavet/sso-oauth/internal/app/ape"
	"github.com/pkg/errors"
)

func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		httpkit.RenderErr(w, problems.BadRequest(errors.New("missing code parameter"))...)
		return
	}

	token, err := h.google.Exchange(r.Context(), code)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	client := h.google.Client(r.Context(), token)
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
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	_, err = h.app.AccountGetByEmail(r.Context(), accountInfo.Email)
	if errors.Is(err, ape.ErrAccountNotFound) {
		session, err := h.app.Login(r.Context(), app.LoginRequest{
			Email: accountInfo.Email,
		})
		if err != nil {
			httpkit.RenderErr(w, problems.InternalError())
			return
		}

		httpkit.Render(w, responses.TokensPair(session.Access, session.Refresh))
		return
	}

	session, err := h.app.Login(r.Context(), app.LoginRequest{
		Email: accountInfo.Email,
	})
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.TokensPair(session.Access, session.Refresh))
}
