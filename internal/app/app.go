package app

import (
	"context"

	"github.com/chains-lab/chains-auth/internal/config"
	"github.com/chains-lab/chains-auth/internal/jwtkit"
	"github.com/chains-lab/chains-auth/internal/repo"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type App struct {
	sessions sessionsRepo
	accounts accountsRepo
	jwt      JWTManager
	log      *logrus.Entry
}

func NewApp(cfg config.Config, log *logrus.Logger) (App, error) {
	sessions, err := repo.NewSessions(cfg, log)
	if err != nil {
		return App{}, err
	}

	accounts, err := repo.NewAccounts(cfg, log)
	if err != nil {
		return App{}, err
	}

	jwt := jwtkit.NewManager(cfg)

	return App{
		sessions: sessions,
		accounts: accounts,
		jwt:      jwt,
		log:      log.WithField("component", "app"),
	}, nil
}

type sessionsRepo interface {
	Create(ctx context.Context, input repo.SessionCreateRequest) error
	Update(ctx context.Context, ID uuid.UUID, input repo.SessionUpdateRequest) error
	Delete(ctx context.Context, ID uuid.UUID) error
	Terminate(ctx context.Context, accountID uuid.UUID) error
	GetByID(ctx context.Context, ID uuid.UUID) (repo.Session, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]repo.Session, error)
	Transaction(fn func(ctx context.Context) error) error
	Drop(ctx context.Context) error
}

type accountsRepo interface {
	Create(ctx context.Context, input repo.AccountCreateRequest) error
	Update(ctx context.Context, ID uuid.UUID, input repo.AccountUpdateRequest) error
	Delete(ctx context.Context, ID uuid.UUID) error
	GetByID(ctx context.Context, ID uuid.UUID) (repo.Account, error)
	GetByEmail(ctx context.Context, email string) (repo.Account, error)
	Transaction(fn func(ctx context.Context) error) error
	Drop(ctx context.Context) error
}

type JWTManager interface {
	EncryptAccess(token string) (string, error)
	EncryptRefresh(token string) (string, error)
	DecryptRefresh(encryptedToken string) (string, error)

	GenerateAccess(
		userID uuid.UUID,
		sessionID uuid.UUID,
		subTypeID uuid.UUID,
		idn roles.Role,
	) (string, error)

	GenerateRefresh(
		userID uuid.UUID,
		sessionID uuid.UUID,
		subTypeID uuid.UUID,
		idn roles.Role,
	) (string, error)
}
