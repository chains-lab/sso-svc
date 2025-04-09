package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/sso-oauth/internal/api/requests"
	"github.com/hs-zavet/sso-oauth/internal/api/responses"
	"github.com/hs-zavet/sso-oauth/internal/app"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewRefresh(r)
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	curToken := req.Data.Attributes.RefreshToken

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		httpkit.RenderErr(w, problems.Unauthorized("Missing Authorization header"))
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		httpkit.RenderErr(w, problems.Unauthorized("Invalid Authorization header format"))
		return
	}

	tokenString := parts[1]

	accountData, err := tokens.VerifyUserJWT(r.Context(), tokenString, h.cfg.JWT.AccessToken.SecretKey)
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		httpkit.RenderErr(w, problems.Unauthorized())
		return
	}

	sessionID := accountData.Session
	//accountRole := accountData.Role
	//subTypeID := accountData.SubID

	accountID, err := uuid.Parse(accountData.Subject)
	if err != nil {
		httpkit.RenderErr(w, problems.Unauthorized("Invalid account ID"))
		return
	}

	//-------------------------------------------------------------------------//

	session, err := h.app.Refresh(r.Context(), accountID, sessionID, app.RefreshRequest{
		Token:  curToken,
		Client: r.Header.Get("User-Agent"),
	})
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.TokensPair(session.Access, session.Refresh))
}
