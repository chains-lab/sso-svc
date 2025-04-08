package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/hs-zavet/sso-oauth/internal/app"
	"github.com/hs-zavet/sso-oauth/internal/app/models"
	"github.com/hs-zavet/sso-oauth/internal/config"
	"github.com/hs-zavet/tokens/roles"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type App interface {
	AccountCreate(ctx context.Context, email string) error
	AccountUpdateSubscription(ctx context.Context, accountID, subscriptionID uuid.UUID) error
	AccountUpdateRole(ctx context.Context, accountID uuid.UUID, role, initiatorRole roles.Role) error
	AccountGetByID(ctx context.Context, accountID uuid.UUID) (models.Account, error)
	AccountGetByEmail(ctx context.Context, email string) (models.Account, error)

	Login(ctx context.Context, request app.LoginRequest) (models.Session, error)
	Refresh(ctx context.Context, accountID, sessionID uuid.UUID, request app.RefreshRequest) (models.Session, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
	TerminateByOwner(ctx context.Context, accountID uuid.UUID) error
	DeleteSessionByOwner(ctx context.Context, sessionID, initiatorSessionID uuid.UUID) error
	TerminateByAdmin(ctx context.Context, userID uuid.UUID) error
	DeleteSessionByAdmin(ctx context.Context, sessionID, initiatorID, initiatorSessionID uuid.UUID) error
	GetSession(ctx context.Context, sessionID uuid.UUID) (models.Session, error)
	GetSessions(ctx context.Context, accountID uuid.UUID) ([]models.Session, error)
}

type Handler struct {
	app    App
	cfg    config.Config
	google oauth2.Config
	log    *logrus.Logger
}

func NewHandlers(app *app.App, cfg config.Config, log *logrus.Logger) *Handler {
	return &Handler{
		app:    app,
		cfg:    cfg,
		google: config.InitGoogleOAuth(cfg),
		log:    log,
	}
}
