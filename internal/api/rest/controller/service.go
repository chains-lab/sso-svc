package controller

import (
	"context"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/chains-lab/sso-svc/internal/domain/services/session"
	"github.com/chains-lab/sso-svc/internal/domain/services/user"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type Session interface {
	Delete(ctx context.Context, sessionID uuid.UUID) error
	DeleteOneForUser(ctx context.Context, userID, sessionID uuid.UUID) error
	DeleteAllForUser(ctx context.Context, userID uuid.UUID) error

	Get(ctx context.Context, sessionID uuid.UUID) (models.Session, error)
	GetForUser(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error)

	ListForUser(
		ctx context.Context,
		userID uuid.UUID,
		page uint,
		size uint,
	) (models.SessionCollection, error)

	Login(ctx context.Context, email, password string) (models.TokensPair, error)
	LoginByGoogle(ctx context.Context, email string) (models.TokensPair, error)

	Refresh(ctx context.Context, oldRefreshToken string) (models.TokensPair, error)
}

type User interface {
	AdminBlockUser(ctx context.Context, userID uuid.UUID) (models.User, error)
	AdminUnblockUser(ctx context.Context, userID uuid.UUID) (models.User, error)

	GetByID(ctx context.Context, ID uuid.UUID) (models.User, error)

	Register(
		ctx context.Context,
		email, pass, role string,
	) (models.User, error)
	RegisterAdmin(
		ctx context.Context,
		initiatorID uuid.UUID,
		email, pass, role string,
	) (models.User, error)

	UpdatePassword(
		ctx context.Context,
		userID uuid.UUID,
		oldPassword, newPassword string,
	) error
}

type Service struct {
	domain domain

	log    logium.Logger
	cfg    internal.Config
	google oauth2.Config
}

type domain struct {
	user    User
	session Session
}

func NewService(cfg internal.Config, log logium.Logger, user user.Service, service session.Service) Service {
	return Service{
		domain: domain{
			user:    user,
			session: service,
		},

		log:    log,
		cfg:    cfg,
		google: cfg.GoogleOAuth(),
	}
}
