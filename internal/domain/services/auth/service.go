package auth

import (
	"context"
	"time"

	"github.com/chains-lab/gatekit/auth"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/google/uuid"
)

type JWTManager interface {
	EncryptAccess(token string) (string, error)
	EncryptRefresh(token string) (string, error)
	DecryptRefresh(encryptedToken string) (string, error)

	ParseRefreshClaims(enc string) (auth.UsersClaims, error)

	GenerateAccess(
		userID uuid.UUID,
		sessionID uuid.UUID,
		idn string,
	) (string, error)

	GenerateRefresh(
		userID uuid.UUID,
		sessionID uuid.UUID,
		role string,
	) (string, error)
}

type Service struct {
	jwt JWTManager
	db  database
}

func New(db database, jwt JWTManager) Service {
	return Service{
		jwt: jwt,
		db:  db,
	}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)

	GetUserPassword(ctx context.Context, userID uuid.UUID) (models.UserPassword, error)

	CreateUser(ctx context.Context, user models.User, pass models.UserPassword) error
	CreateSession(ctx context.Context, session models.Session, token string) error

	GetSession(ctx context.Context, sessionID uuid.UUID) (models.Session, error)
	GetSessionToken(ctx context.Context, sessionID uuid.UUID) (string, error)

	UpdateSessionToken(ctx context.Context, sessionID uuid.UUID, token string, updatedAt time.Time) error
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, passwordHash string, updatedAt time.Time) error

	DeleteAllSessionsForUser(ctx context.Context, userID uuid.UUID) error
}
