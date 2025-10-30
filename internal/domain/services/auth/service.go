package auth

import (
	"context"
	"time"

	"github.com/chains-lab/restkit/token"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/google/uuid"
)

type JWTManager interface {
	EncryptAccess(token string) (string, error)
	EncryptRefresh(token string) (string, error)
	DecryptRefresh(encryptedToken string) (string, error)

	ParseRefreshClaims(enc string) (token.UsersClaims, error)

	GenerateAccess(
		user models.User, sessionID uuid.UUID,
	) (string, error)

	GenerateRefresh(
		user models.User, sessionID uuid.UUID,
	) (string, error)
}

type passwordManager interface {
	ReliabilityCheck(password string) error
	CheckPasswordMatch(password, hash string) error
}

type Service struct {
	jwt  JWTManager
	db   database
	pass passwordManager
}

func New(db database, jwt JWTManager, pass passwordManager) Service {
	return Service{
		jwt:  jwt,
		db:   db,
		pass: pass,
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

	UpdateUserCity(
		ctx context.Context,
		userID uuid.UUID,
		cityID *uuid.UUID,
		cityRole *string,
		updatedAt time.Time,
	) error

	UpdateUserCompany(
		ctx context.Context,
		userID uuid.UUID,
		companyID *uuid.UUID,
		companyRole *string,
		updatedAt time.Time,
	) error

	DeleteAllSessionsForUser(ctx context.Context, userID uuid.UUID) error
}
