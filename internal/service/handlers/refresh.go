package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/recovery-flow/comtools/cifractx"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/sectools"
	"github.com/recovery-flow/sso-oauth/internal/service/requests"
	"github.com/recovery-flow/sso-oauth/internal/service/responses"
	"github.com/sirupsen/logrus"
)

func Refresh(w http.ResponseWriter, r *http.Request) {
	server, err := cifractx.GetValue[*config.Server](r.Context(), config.SERVER)
	if err != nil {
		logrus.Errorf("Failed to retrieve service configuration %s", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}
	log := server.Logger

	req, err := requests.NewRefresh(r)
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	refreshToken := req.Data.Attributes.RefreshToken

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Debugf("Missing Authorization header")
		httpkit.RenderErr(w, problems.Unauthorized("Missing Authorization header"))
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		log.Debugf("Invalid Authorization header format")
		httpkit.RenderErr(w, problems.Unauthorized("Invalid Authorization header format"))
		return
	}
	tokenString := parts[1]

	userData, err := server.TokenManager.VerifyJWTAndExtractClaims(tokenString, server.Config.JWT.AccessToken.SecretKey)
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		log.Warnf("Token validation failed: %v", err)
		httpkit.RenderErr(w, problems.Unauthorized())
		return
	}

	userID := userData.ID
	sessionID := userData.DevID

	user, err := server.SqlDB.Accounts.GetById(r, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.Unauthorized())
			return
		}
		log.Errorf("Failed to get user: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	session, err := server.SqlDB.Sessions.GetByID(r, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.Unauthorized("session not found"))
			return
		}
		log.Errorf("Failed to get session: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	log.Debugf("Session token: %s \n EncryptionKey: %s`", session.Token, server.Config.JWT.RefreshToken.EncryptionKey)

	decryptedToken, err := sectools.DecryptToken(session.Token, server.Config.JWT.RefreshToken.EncryptionKey)
	if err != nil {
		log.Errorf("Failed to decrypt refresh token: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	if decryptedToken != refreshToken {
		server.Logger.Warn("Provided refresh token does not match the stored token")
		httpkit.RenderErr(w, problems.Conflict())
		return
	}

	tokenAccess, err := server.TokenManager.GenerateJWT(user.ID, sessionID, user.Role, server.Config.JWT.AccessToken.TokenLifetime, server.Config.JWT.AccessToken.SecretKey)
	if err != nil {
		server.Logger.Errorf("Error generating access token: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	tokenRefresh, err := server.TokenManager.GenerateJWT(user.ID, sessionID, user.Role, server.Config.JWT.RefreshToken.TokenLifetime, server.Config.JWT.RefreshToken.SecretKey)
	if err != nil {
		server.Logger.Errorf("Error generating refresh token: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	encryptedToken, err := sectools.EncryptToken(tokenRefresh, server.Config.JWT.RefreshToken.EncryptionKey)
	if err != nil {
		log.Errorf("Failed to encrypt refresh token: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	err = server.SqlDB.Sessions.UpdateToken(r, userID, encryptedToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.Unauthorized())
			return
		}
		log.Errorf("Error updating last used and refresh token: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.TokensPair(tokenAccess, tokenRefresh))
}
