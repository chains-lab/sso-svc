package controller

import (
	"context"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type SessionSvc interface {
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
	) (models.SessionsCollection, error)
}

type UserSvc interface {
	GetByID(ctx context.Context, ID uuid.UUID) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
}

type AuthSvc interface {
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

	Login(ctx context.Context, email, password string) (models.TokensPair, error)
	LoginByGoogle(ctx context.Context, email string) (models.TokensPair, error)
	CreateSession(ctx context.Context, userID uuid.UUID, role string) (models.TokensPair, error)

	Refresh(ctx context.Context, oldRefreshToken string) (models.TokensPair, error)
}

type services struct {
	Session SessionSvc
	User    UserSvc
	Auth    AuthSvc
}

type Service struct {
	google oauth2.Config
	domain services
	log    logium.Logger
}

func NewService(log logium.Logger, google oauth2.Config, user UserSvc, session SessionSvc, auth AuthSvc) *Service {
	return &Service{
		log:    log,
		google: google,
		domain: services{
			Session: session,
			User:    user,
			Auth:    auth,
		},
	}
}
