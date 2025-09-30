package data

import (
	"context"
	"time"

	"github.com/chains-lab/sso-svc/internal/data/schemas"
	"github.com/google/uuid"
)

func (d *Database) CreateSession(ctx context.Context, session schemas.Session) error {
	return d.sql.sessions.New().Insert(ctx, session)
}

func (d *Database) GetSession(ctx context.Context, sessionID uuid.UUID) (schemas.Session, error) {
	return d.sql.sessions.New().FilterID(sessionID).Get(ctx)
}

func (d *Database) GetOneSessionForUser(ctx context.Context, userID, sessionID uuid.UUID) (schemas.Session, error) {
	return d.sql.sessions.New().FilterUserID(userID).FilterID(sessionID).Get(ctx)
}

func (d *Database) GetAllSessionsForUser(
	ctx context.Context,
	userID uuid.UUID,
	limit, offset uint,
) ([]schemas.Session, uint, error) {
	sessions, err := d.sql.sessions.New().FilterUserID(userID).Page(limit, offset).OrderCreatedAt(false).Select(ctx)
	if err != nil {
		return nil, 0, err
	}

	count, err := d.sql.sessions.New().FilterUserID(userID).Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return sessions, count, nil
}

func (d *Database) UpdateSessionToken(
	ctx context.Context,
	sessionID uuid.UUID,
	token string,
	lastUsedAt time.Time,
) error {
	return d.sql.sessions.New().FilterID(sessionID).UpdateToken(token).UpdateLastUsed(lastUsedAt).Update(ctx)
}

func (d *Database) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	return d.sql.sessions.New().FilterID(sessionID).Delete(ctx)
}

func (d *Database) DeleteOneSessionForUser(ctx context.Context, userID, sessionID uuid.UUID) error {
	return d.sql.sessions.New().FilterUserID(userID).FilterID(sessionID).Delete(ctx)
}

func (d *Database) DeleteAllSessionsForUser(ctx context.Context, userID uuid.UUID) error {
	return d.sql.sessions.New().FilterUserID(userID).Delete(ctx)
}
