package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/recovery-flow/comtools/cifractx"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/comtools/jsonkit"
	"github.com/recovery-flow/roles"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/sectools"
	"github.com/recovery-flow/sso-oauth/internal/service/events"
	"github.com/recovery-flow/sso-oauth/internal/service/responses"
	"github.com/recovery-flow/sso-oauth/internal/service/utils"
)

func LogSimple(w http.ResponseWriter, r *http.Request) {
	server, err := cifractx.GetValue[*config.Server](r.Context(), config.SERVER)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}
	log := server.Logger

	if !server.Config.Email.Off {
		log.Info("Email is on")
		httpkit.RenderErr(w, problems.Forbidden("Email is on"))
	}

	type emailReq struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	var req emailReq

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = jsonkit.NewDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"email": validation.Validate(req.Email, validation.Required),
		"role":  validation.Validate(req.Role, validation.Required),
	}
	if errs.Filter() != nil {
		log.WithError(err).Error("Failed to parse email")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	account, err := server.SqlDB.Accounts.GetByEmail(r, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			role, err := roles.StringToRoleUser(req.Role)
			if err != nil {
				log.Errorf("error getting role: %v", err)
				httpkit.RenderErr(w, problems.InternalError())
				return
			}

			account, err = server.SqlDB.Accounts.Create(r, req.Email, role)
			if err != nil {
				log.Errorf("error creating user: %v", err)
				httpkit.RenderErr(w, problems.InternalError())
				return
			}

			event := events.AccountCreated{
				Event:     "AccountCreate",
				UserID:    account.ID.String(),
				Email:     req.Email,
				Role:      string(role),
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
