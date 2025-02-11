package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/rerabbit"
	"github.com/recovery-flow/sso-oauth/internal/sectools"
	"github.com/recovery-flow/sso-oauth/internal/service/events/entities"
	"github.com/recovery-flow/sso-oauth/internal/service/responses"
)

func (h *Handlers) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	svc := h.svc
	log := svc.Logger

	code := r.URL.Query().Get("code")
	if code == "" {
		log.Debugf("missing code parameter")
		httpkit.RenderErr(w, problems.BadRequest(errors.New("missing code parameter"))...)
		return
	}

	token, err := svc.GoogleOAuth.Exchange(r.Context(), code)
	if err != nil {
		log.Errorf("failed to exchange code for token: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	client := svc.GoogleOAuth.Client(r.Context(), token)
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

	account, err := svc.SqlDB.Accounts.GetByEmail(r, userInfo.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			account, err = svc.SqlDB.Accounts.Create(r, userInfo.Email, "user")
			if err != nil {
				log.Errorf("error creating user: %v", err)
				httpkit.RenderErr(w, problems.InternalError())
				return
			}
			event := entities.AccountCreated{
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
			err = svc.Rabbit.PublishJSON(r.Context(), body, rerabbit.PublishOptions{
				Exchange:   "re-news.sso",
				RoutingKey: "account.created",
			})
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

	tokenAccess, tokenRefresh, err := sectools.GenerateTokens(*svc, account, deviceID)
	if err != nil {
		log.Errorf("error generating tokens: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	tokenCrypto, err := sectools.EncryptToken(tokenRefresh, svc.Config.JWT.RefreshToken.EncryptionKey)
	if err != nil {
		log.Errorf("error encrypting token: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	_, err = svc.SqlDB.Sessions.Create(r, account.ID, deviceID, tokenCrypto)
	if err != nil {
		log.Errorf("error creating session: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.TokensPair(tokenAccess, tokenRefresh))
}
