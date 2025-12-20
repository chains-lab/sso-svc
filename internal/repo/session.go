package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/umisto/restkit/pagi"
	"github.com/umisto/sso-svc/internal/domain/entity"
	"github.com/umisto/sso-svc/internal/repo/pgdb"
)

func (r *Repository) CreateSession(ctx context.Context, sessionID, accountID uuid.UUID, hashToken string) (entity.Session, error) {
	now := time.Now().UTC()

	row := pgdb.Session{
		ID:        sessionID,
		AccountID: accountID,
		HashToken: hashToken,
		LastUsed:  now,
		CreatedAt: now,
	}

	err := r.sql.sessions.Insert(ctx, row)
	if err != nil {
		return entity.Session{}, err
	}

	return row.ToEntity(), nil
}

func (r *Repository) GetSession(ctx context.Context, sessionID uuid.UUID) (entity.Session, error) {
	row, err := r.sql.sessions.New().FilterID(sessionID).Get(ctx)
	switch {
	case err != nil:
		return entity.Session{}, err
	case row.ID == uuid.Nil:
		return entity.Session{}, nil
	}

	return row.ToEntity(), nil
}

func (r *Repository) GetAccountSession(ctx context.Context, userID, sessionID uuid.UUID) (entity.Session, error) {
	row, err := r.sql.sessions.New().
		FilterID(sessionID).
		FilterAccountID(userID).
		Get(ctx)
	switch {
	case err != nil:
		return entity.Session{}, err
	case row.ID == uuid.Nil:
		return entity.Session{}, nil
	}

	return row.ToEntity(), nil
}

func (r *Repository) GetSessionsForAccount(ctx context.Context, userID uuid.UUID, page, size int32) (entity.SessionsCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	rows, err := r.sql.sessions.New().
		FilterAccountID(userID).
		OrderCreatedAt(false).
		Page(uint64(limit), uint64(offset)).
		Select(ctx)
	if err != nil {
		return entity.SessionsCollection{}, err
	}

	total, err := r.sql.sessions.New().
		FilterAccountID(userID).
		Count(ctx)
	if err != nil {
		return entity.SessionsCollection{}, err
	}

	result := make([]entity.Session, 0, len(rows))
	for _, s := range rows {
		result = append(result, s.ToEntity())
	}

	return entity.SessionsCollection{
		Data:  result,
		Page:  page,
		Size:  size,
		Total: int64(total),
	}, nil
}

func (r *Repository) GetSessionToken(ctx context.Context, sessionID uuid.UUID) (string, error) {
	row, err := r.sql.sessions.New().FilterID(sessionID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return "", nil
	case err != nil:
		return "", err
	}

	return row.HashToken, nil
}

func (r *Repository) UpdateSessionToken(ctx context.Context, sessionID uuid.UUID, token string) (entity.Session, error) {
	sess, err := r.sql.sessions.New().
		FilterID(sessionID).
		UpdateToken(token).
		Update(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.Session{}, nil
	case err != nil:
		return entity.Session{}, err
	}

	if len(sess) != 1 {
		return entity.Session{}, fmt.Errorf("expected 1 session, got %d", len(sess))
	}
	return sess[0].ToEntity(), nil
}

func (r *Repository) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	return r.sql.sessions.New().FilterID(sessionID).Delete(ctx)
}

func (r *Repository) DeleteSessionsForAccount(ctx context.Context, userID uuid.UUID) error {
	return r.sql.sessions.New().FilterAccountID(userID).Delete(ctx)
}

func (r *Repository) DeleteAccountSession(ctx context.Context, userID, sessionID uuid.UUID) error {
	return r.sql.sessions.New().
		FilterID(sessionID).
		FilterAccountID(userID).
		Delete(ctx)
}

func toSessionModel(s pgdb.Session) entity.Session {
	return entity.Session{
		ID:        s.ID,
		AccountID: s.AccountID,
		CreatedAt: s.CreatedAt,
		LastUsed:  s.LastUsed,
	}
}
