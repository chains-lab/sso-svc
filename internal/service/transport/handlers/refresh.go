package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/core/tools"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/requests"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/responses"
)

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewRefresh(r)
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	refreshToken := req.Data.Attributes.RefreshToken

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.Log.Debugf("Missing Authorization header")
		httpkit.RenderErr(w, problems.Unauthorized("Missing Authorization header"))
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		h.Log.Debugf("Invalid Authorization header format")
		httpkit.RenderErr(w, problems.Unauthorized("Invalid Authorization header format"))
		return
	}

	tokenString := parts[1]

	h.Log.Debugf("Token received: %s", tokenString)

	userData, err := h.svc.TokenManager.VerifyJWT(r.Context(), tokenString)
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		h.Log.Warnf("Token validation failed: %v", err)
		httpkit.RenderErr(w, problems.Unauthorized())
		return
	}
	if userData == nil {
		h.Log.Debugf("Token validation failed")
		httpkit.RenderErr(w, problems.Unauthorized("Token validation failed"))
		return
	}

	userID, err := uuid.Parse(userData.ID)
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	sessionID, err := uuid.Parse(*userData.SessionID)
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	user, err := h.Domain.Account.Get(r.Context(), userID)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	session, err := h.Domain.Session.Get(r.Context(), sessionID)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	err = h.Domain.Session.TokenRefValidate(r.Context(), session, h.svc.Config.JWT.RefreshToken.EncryptionKey, refreshToken)
	if err != nil {
		httpkit.RenderErr(w, problems.Unauthorized("Invalid refresh token"))
		return
	}

	sesIDStr := sessionID.String()
	tokenAccess, err := svc.TokenManager.GenerateJWT(
		svc.Config.Server.Name,
		userID.String(),
		svc.Config.JWT.AccessToken.TokenLifetime,
		nil,
		&user.Role,
		&sesIDStr,
	)
	if err != nil {
		svc.Logger.Errorf("Error generating access token: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	tokenRefresh, err := svc.TokenManager.GenerateJWT(
		svc.Config.Server.Name,
		userID.String(),
		svc.Config.JWT.RefreshToken.TokenLifetime,
		nil,
		&user.Role,
		&sesIDStr,
	)
	if err != nil {
		svc.Logger.Errorf("Error generating refresh token: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	encryptedToken, err := tools.EncryptToken(tokenRefresh, svc.Config.JWT.RefreshToken.EncryptionKey)
	if err != nil {
		log.Errorf("Failed to encrypt refresh token: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	_, err = svc.DB.Sessions.UpdateToken(r, userID, encryptedToken)
	if err != nil {
		render.RenderSelectErr(w, log, err, "Failed to update session token")
		return
	}

	httpkit.Render(w, responses.TokensPair(tokenAccess, tokenRefresh))
}
