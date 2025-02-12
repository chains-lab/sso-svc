package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/comtools/jsonkit"
	"github.com/recovery-flow/rerabbit"
	"github.com/recovery-flow/roles"
	"github.com/recovery-flow/sso-oauth/internal/sectools"
	"github.com/recovery-flow/sso-oauth/internal/service/events/entities"
	"github.com/recovery-flow/sso-oauth/internal/service/responses"
)

func (h *Handlers) LoginSimple(w http.ResponseWriter, r *http.Request) {
	svc := h.svc
	log := svc.Logger

	if !svc.Config.Email.Off {
		log.Info("Email is on")
		httpkit.RenderErr(w, problems.Forbidden("Email is on"))
	}

	type emailReq struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	var req emailReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = jsonkit.NewDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"email": validation.Validate(req.Email, validation.Required),
		"role":  validation.Validate(req.Role, validation.Required),
	}
	if errs.Filter() != nil {
		log.WithError(errs.Filter()).Error("Failed to parse email")
		httpkit.RenderErr(w, problems.BadRequest(errs.Filter())...)
		return
	}

	account, err := svc.DB.Accounts.GetByEmail(r, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			role, err := roles.StringToRoleUser(req.Role)
			if err != nil {
				log.Errorf("error getting role: %v", err)
				httpkit.RenderErr(w, problems.InternalError())
				return
			}

			account, err = svc.DB.Accounts.Create(r, req.Email, role)
			if err != nil {
				log.Errorf("error creating user: %v", err)
				httpkit.RenderErr(w, problems.InternalError())
				return
			}

			event := entities.AccountCreated{
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
			err = svc.Rabbit.PublishWithRetry(r.Context(), rerabbit.PublishOptions{
				Exchange:     "re-news.sso",
				RoutingKey:   "account.created",
				Mandatory:    true,
				Body:         body,
				DeliveryMode: 2,
			}, 3, 2*time.Second)
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

	tokenAccess, tokenRefresh, err := sectools.GenerateTokens(*svc, *account, deviceID)
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

	_, err = svc.DB.Sessions.Create(r, account.ID, deviceID, tokenCrypto)
	if err != nil {
		log.Errorf("error creating session: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.TokensPair(tokenAccess, tokenRefresh))
}
