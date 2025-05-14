package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/chains-lab/chains-auth/internal/api/requests"
	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/chains-auth/internal/app"
	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/google/uuid"
)

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewRefresh(r)
	if err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Error:  err,
		})...)
		return
	}
	curToken := req.Data.Attributes.RefreshToken

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusUnauthorized,
			Detail: "Missing Authorization header",
		})...)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusUnauthorized,
			Detail: "Invalid Authorization header",
		})...)
		return
	}

	tokenString := parts[1]

	userData, err := tokens.VerifyAccountsJWT(r.Context(), tokenString, h.cfg.JWT.AccessToken.SecretKey)
	if err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusUnauthorized,
			Title:  "Unauthorized",
			Detail: "Token validation failed",
		})...)
		return
	}

	sessionID := userData.Session
	//accountRole := userData.Role
	//subTypeID := userData.SubID

	accountID, err := uuid.Parse(userData.Subject)
	if err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Detail: "Account ID must be a valid UUID.",
		})...)
		return
	}

	//-------------------------------------------------------------------------//

	session, err := h.app.Refresh(r.Context(), accountID, sessionID, app.RefreshRequest{
		Token:  curToken,
		Client: r.Header.Get("User-Agent"),
	})
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrSessionDoseNotExits):
			h.log.WithError(err).Errorf("session not found session id: %s", sessionID)
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusNotFound,
				Title:  "Session not found",
				Detail: "Session does not exist.",
			})...)
		case errors.Is(err, ape.ErrSessionsClientMismatch):
			h.log.WithError(err).Errorf("session client mismatch session id: %s", sessionID)
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusConflict,
				Detail: "Session client does not match.",
			})...)
		case errors.Is(err, ape.ErrSessionsTokenMismatch):
			h.log.WithError(err).Errorf("session token mismatch session id: %s", sessionID)
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusConflict,
				Detail: "Session token mismatch.",
			})...)
		default:
			h.log.WithError(err).Error("error refreshing session")
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}
		return
	}

	httpkit.Render(w, responses.TokensPair(session.Access, session.Refresh))
}
