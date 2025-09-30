package auth

import (
	"context"
	"time"

	"github.com/chains-lab/gatekit/auth"
	"github.com/chains-lab/sso-svc/internal/data/schemas"
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

	GetUserByID(ctx context.Context, userID uuid.UUID) (schemas.User, error)
	GetUserByEmail(ctx context.Context, email string) (schemas.User, error)

	CreateUser(ctx context.Context, user schemas.User) error
	CreateSession(ctx context.Context, session schemas.Session) error

	GetSession(ctx context.Context, sessionID uuid.UUID) (schemas.Session, error)

	UpdateSessionToken(ctx context.Context, sessionID uuid.UUID, token string, updatedAt time.Time) error
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, passwordHash string, updatedAt time.Time) error

	DeleteAllSessionsForUser(ctx context.Context, userID uuid.UUID) error
}
