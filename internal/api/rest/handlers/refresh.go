package handlers

import (
	"net/http"
	"strings"

	"github.com/chains-lab/chains-auth/internal/api/rest/requests"
	"github.com/chains-lab/chains-auth/internal/api/rest/responses"
	"github.com/chains-lab/chains-auth/internal/app"
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
			Code:   "MISSING_AUTHORIZATION_HEADER",
			Detail: "Missing Authorization header",
		})...)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusUnauthorized,
			Code:   "MISSING_AUTHORIZATION_HEADER",
			Detail: "Missing Authorization header",
		})...)
		return
	}

	tokenString := parts[1]

	userData, err := tokens.VerifyAccountsJWT(r.Context(), tokenString, h.cfg.JWT.AccessToken.SecretKey)
	if err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusUnauthorized,
			Code:   "TOKEN_VALIDATION_FAILED",
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

	requestID := uuid.New()
	log := h.log.WithField("request_id", requestID)

	session, appErr := h.app.Refresh(r.Context(), accountID, sessionID, app.RefreshRequest{
		Token:  curToken,
		Client: r.Header.Get("User-Agent"),
	})
	if appErr != nil {
		h.controllers.ResultFromApp(w, requestID, appErr)
		return
	}

	log.Infof("Session %s refreshed successfully", session.ID)
	httpkit.Render(w, responses.TokensPair(session.Access, session.Refresh))
}
