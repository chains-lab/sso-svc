package session

import (
	"context"
	"database/sql"

	"github.com/chains-lab/gatekit/auth"
	"github.com/chains-lab/sso-svc/internal/app/jwtmanager"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/google/uuid"
)

type sessionsQ interface {
	New() dbx.SessionsQ
	Insert(ctx context.Context, input dbx.Session) error
	Update(ctx context.Context, input map[string]any) error
	Delete(ctx context.Context) error
	Select(ctx context.Context) ([]dbx.Session, error)
	Get(ctx context.Context) (dbx.Session, error)

	FilterID(id uuid.UUID) dbx.SessionsQ
	FilterUserID(userID uuid.UUID) dbx.SessionsQ

	Page(limit, offset uint64) dbx.SessionsQ
	Count(ctx context.Context) (uint64, error)

	OrderCreatedAt(ascending bool) dbx.SessionsQ
}

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
		emailVerified bool,
	) (string, error)
}

type Session struct {
	query sessionsQ

	jwt JWTManager
}

func CreateSession(pg *sql.DB, manager jwtmanager.Manager) Session {
	return Session{
		query: dbx.NewSessions(pg),
		jwt:   manager,
	}
}
