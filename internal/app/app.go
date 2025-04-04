package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/hs-zavet/sso-oauth/internal/config"
	"github.com/hs-zavet/sso-oauth/internal/jwtkit"
	"github.com/hs-zavet/sso-oauth/internal/repo"
	"github.com/hs-zavet/tokens/identity"
)

type App struct {
	sessions sessionsRepo
	accounts accountsRepo
	jwt      JWTManager
}

func NewApp(cfg *config.Config) (App, error) {
	sessions, err := repo.NewSessions(cfg)
	if err != nil {
		return App{}, err
	}

	accounts, err := repo.NewAccounts(cfg)
	if err != nil {
		return App{}, err
	}

	jwt := jwtkit.NewManager(cfg)

	return App{
		sessions: sessions,
		accounts: accounts,
		jwt:      jwt,
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
}

type accountsRepo interface {
	Create(ctx context.Context, input repo.AccountCreateRequest) error
	Update(ctx context.Context, ID uuid.UUID, input repo.AccountUpdateRequest) error
	Delete(ctx context.Context, ID uuid.UUID) error
	GetByID(ctx context.Context, ID uuid.UUID) (repo.Account, error)
	GetByEmail(ctx context.Context, email string) (repo.Account, error)
	Transaction(fn func(ctx context.Context) error) error
}

type JWTManager interface {
	EncryptAccess(token string) (string, error)
	EncryptRefresh(token string) (string, error)
	DecryptRefresh(encryptedToken string) (string, error)

	GenerateAccess(
		userID uuid.UUID,
		sessionID uuid.UUID,
		subTypeID uuid.UUID,
		idn identity.IdnType,
	) (string, error)

	GenerateRefresh(
		userID uuid.UUID,
		sessionID uuid.UUID,
		subTypeID uuid.UUID,
		idn identity.IdnType,
	) (string, error)
}
