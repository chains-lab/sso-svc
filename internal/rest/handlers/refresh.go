package handlers

import (
	"net/http"
	"strings"

	"github.com/chains-lab/chains-auth/internal/rest/requests"
	"github.com/chains-lab/chains-auth/internal/rest/responses"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/google/uuid"
)

func (h *Handlers) Refresh(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()

	//------------------------------- COPY FROM MDLV ------------------------------------------//

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

	userData, err := tokens.VerifyUserJWT(r.Context(), tokenString, h.cfg.JWT.AccessToken.SecretKey)
	if err != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusUnauthorized,
			Code:   "TOKEN_VALIDATION_FAILED",
			Detail: "Token validation failed",
		})...)
		return
	}

	sessionID := userData.Session

	userID, err := uuid.Parse(userData.Subject)
	if err != nil {
		h.presenter.InvalidParameter(w, uuid.Nil, err, "user_id")
		return
	}

	//-------------------------------END COPY FROM MDLV------------------------------------------//

	req, err := requests.NewRefresh(r)
	if err != nil {
		h.presenter.InvalidPointer(w, requestID, err)
		return
	}

	curToken := req.Data.Attributes.RefreshToken

	log := h.log.WithField("request_id", requestID)

	session, tokensPair, appErr := h.app.Refresh(r.Context(), userID, sessionID, r.Header.Get("User-Agent"), curToken)
	if appErr != nil {
		h.presenter.AppError(w, requestID, appErr)
		return
	}

	log.Infof("Session %s refreshed successfully", session.ID)
	httpkit.Render(w, responses.TokensPair(tokensPair.Access, tokensPair.Refresh))
}
