package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/recovery-flow/comtools/cifractx"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/sectools"
	"github.com/recovery-flow/sso-oauth/internal/service/events"
	"github.com/recovery-flow/sso-oauth/internal/service/responses"
	"github.com/recovery-flow/sso-oauth/internal/service/utils"
)

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	server, err := cifractx.GetValue[*config.Server](r.Context(), config.SERVER)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}
	log := server.Logger

	code := r.URL.Query().Get("code")
	if code == "" {
		log.Debugf("missing code parameter")
		httpkit.RenderErr(w, problems.BadRequest(errors.New("missing code parameter"))...)
		return
	}

	token, err := server.GoogleOAuth.Exchange(r.Context(), code)
	if err != nil {
		log.Errorf("failed to exchange code for token: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	client := server.GoogleOAuth.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Errorf("failed to get user info: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Errorf("failed to decode user info: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	account, err := server.SqlDB.Accounts.GetByEmail(r, userInfo.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			account, err = server.SqlDB.Accounts.Create(r, userInfo.Email, "user")
			if err != nil {
				log.Errorf("error creating user: %v", err)
				httpkit.RenderErr(w, problems.InternalError())
				return
			}
			event := events.AccountCreated{
				Event:     "AccountCreate",
				UserID:    account.ID.String(),
				Email:     userInfo.Email,
				Role:      "user",
				Timestamp: time.Now().UTC(),
			}

			body, err := json.Marshal(event)
			if err != nil {
				log.Errorf("error serializing event: %v", err)
				httpkit.RenderErr(w, problems.InternalError())
				return
			}
			err = server.Broker.Publish(
				server.Config.Rabbit.Exchange,
				"account",
				"account.create",
				body)
			if err != nil {
				log.Errorf("error publishing event: %v", err)
				httpkit.RenderErr(w, problems.InternalError())
				return
			}
		} else {
			log.Errorf("error getting user: %v", err)
			httpkit.RenderErr(w, problems.InternalError())
			return
		}
	}

	deviceID := uuid.New()

	tokenAccess, tokenRefresh, err := utils.GenerateTokens(*server, account, deviceID)
	if err != nil {
		log.Errorf("error generating tokens: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	tokenCrypto, err := sectools.EncryptToken(tokenRefresh, server.Config.JWT.RefreshToken.EncryptionKey)
	if err != nil {
		log.Errorf("error encrypting token: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	_, err = server.SqlDB.Sessions.Create(r, account.ID, deviceID, tokenCrypto)
	if err != nil {
		log.Errorf("error creating session: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.TokensPair(tokenAccess, tokenRefresh))
}
