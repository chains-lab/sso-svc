package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/data/pgdb"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (d *Database) CreateSession(ctx context.Context, session models.Session, token string) error {
	return d.sql.sessions.New().Insert(ctx, pgdb.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		Token:     token,
		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	})
}

func (d *Database) GetSession(ctx context.Context, sessionID uuid.UUID) (models.Session, error) {
	row, err := d.sql.sessions.New().FilterID(sessionID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Session{}, nil
	case err != nil:
		return models.Session{}, err
	}

	return sessionSchemaToModel(row), nil
}

func (d *Database) GetOneSessionForUser(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error) {
	row, err := d.sql.sessions.New().FilterUserID(userID).FilterID(sessionID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Session{}, nil
	case err != nil:
		return models.Session{}, err
	}

	return sessionSchemaToModel(row), nil
}

func (d *Database) GetAllSessionsForUser(
	ctx context.Context,
	userID uuid.UUID,
	page, size uint64,
) (models.SessionsCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	sessions, err := d.sql.sessions.New().FilterUserID(userID).Page(limit, offset).OrderCreatedAt(false).Select(ctx)
	if err != nil {
		return models.SessionsCollection{}, err
	}

	total, err := d.sql.sessions.New().FilterUserID(userID).Count(ctx)
	if err != nil {
		return models.SessionsCollection{}, err
	}

	result := make([]models.Session, 0, len(sessions))
	for _, s := range sessions {
		result = append(result, sessionSchemaToModel(s))
	}

	return models.SessionsCollection{
		Data:  result,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (d *Database) GetSessionToken(ctx context.Context, sessionID uuid.UUID) (string, error) {
	row, err := d.sql.sessions.New().FilterID(sessionID).Get(ctx)
	if err != nil {
		return "", err
	}

	return row.Token, nil
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

func sessionSchemaToModel(s pgdb.Session) models.Session {
	return models.Session{
		ID:        s.ID,
		UserID:    s.UserID,
		LastUsed:  s.LastUsed,
		CreatedAt: s.CreatedAt,
	}
}

func sessionModelToSchema(s models.Session) pgdb.Session {
	return pgdb.Session{
		ID:        s.ID,
		UserID:    s.UserID,
		LastUsed:  s.LastUsed,
		CreatedAt: s.CreatedAt,
	}
}
