package controller

import (
	"context"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal/domain"
	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type core interface {
	Registration(
		ctx context.Context,
		params domain.RegistrationParams,
	) (entity.Account, error)
	RegistrationByAdmin(
		ctx context.Context,
		initiatorID uuid.UUID,
		params domain.RegistrationParams,
	) (entity.Account, error)

	LoginByEmail(ctx context.Context, email, password string) (entity.TokensPair, error)
	LoginByUsername(ctx context.Context, username, password string) (entity.TokensPair, error)
	LoginByGoogle(ctx context.Context, email string) (entity.TokensPair, error)

	Refresh(ctx context.Context, oldRefreshToken string) (entity.TokensPair, error)

	UpdatePassword(
		ctx context.Context,
		accountID uuid.UUID,
		oldPassword, newPassword string,
	) error
	UpdateUsername(
		ctx context.Context,
		accountID uuid.UUID,
		password string,
		newUsername string,
	) (entity.Account, error)

	GetAccountByID(ctx context.Context, ID uuid.UUID) (entity.Account, error)
	GetSessionForAccount(ctx context.Context, accountID, sessionID uuid.UUID) (entity.Session, error)
	GetSessionsForAccount(
		ctx context.Context,
		accountID uuid.UUID,
		page int32,
		size int32,
	) (entity.SessionsCollection, error)

	GetAccountEmailData(ctx context.Context, ID uuid.UUID) (entity.AccountEmail, error)

	DeleteOwnAccount(ctx context.Context, accountID uuid.UUID) error
	DeleteOwnSession(ctx context.Context, accountID, sessionID uuid.UUID) error
	DeleteOwnSessions(ctx context.Context, accountID uuid.UUID) error
	Logout(ctx context.Context, accountID, sessionID uuid.UUID) error
}

type Service struct {
	google oauth2.Config
	domain core
	log    logium.Logger
}

func New(log logium.Logger, google oauth2.Config, domain core) *Service {
	return &Service{
		log:    log,
		google: google,
		domain: domain,
	}
}
