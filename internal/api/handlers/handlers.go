package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/hs-zavet/sso-oauth/internal/app"
	"github.com/hs-zavet/sso-oauth/internal/app/models"
	"github.com/hs-zavet/sso-oauth/internal/config"
	"github.com/hs-zavet/sso-oauth/internal/repo"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type App interface {
	AccountCreate(ctx context.Context, email string) error
	AccountUpdateSubscription(ctx context.Context, ID uuid.UUID, subscriptionID uuid.UUID) error
	AccountUpdateRole(ctx context.Context, ID uuid.UUID, role string) error
	AccountGetByID(ctx context.Context, ID uuid.UUID) (repo.Account, error)
	AccountGetByEmail(ctx context.Context, email string) (repo.Account, error)

	Login(ctx context.Context, request app.LoginRequest) (models.Session, error)
	Refresh(ctx context.Context, sessionID uuid.UUID, request app.RefreshRequest) (models.Session, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
	Terminate(ctx context.Context, sessionID uuid.UUID) error
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	GetSession(ctx context.Context, sessionID uuid.UUID) (models.Session, error)
	GetSessions(ctx context.Context, accountID uuid.UUID) ([]models.Session, error)
}

type Handler struct {
	log    *logrus.Logger
	google oauth2.Config
	app    App
	cfg    *config.Config
}

func NewHandlers(cfg *config.Config, app *app.App) *Handler {
	return &Handler{
		log:    cfg.Log,
		app:    app,
		google: config.InitGoogleOAuth(cfg),
		cfg:    cfg,
	}
}
