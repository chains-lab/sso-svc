package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/chains-lab/restkit/pagi"
	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/repo/pgdb"
	"github.com/google/uuid"
)

func (r *Repository) CreateSession(ctx context.Context, sessionID, accountID uuid.UUID, hashToken string) (entity.Session, error) {
	res, err := r.sql.CreateSession(ctx, pgdb.CreateSessionParams{
		ID:        sessionID,
		AccountID: accountID,
		HashToken: hashToken,
	})
	if err != nil {
		return entity.Session{}, err
	}

	return res.ToEntity(), nil
}

func (r *Repository) GetSession(ctx context.Context, sessionID uuid.UUID) (entity.Session, error) {
	res, err := r.sql.GetSessionByID(ctx, sessionID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.Session{}, nil
	case err != nil:
		return entity.Session{}, err
	}

	return res.ToEntity(), nil
}

func (r *Repository) GetAccountSession(
	ctx context.Context,
	userID, sessionID uuid.UUID,
) (entity.Session, error) {
	res, err := r.sql.GetAccountSession(ctx, pgdb.GetAccountSessionParams{
		ID:        sessionID,
		AccountID: userID,
	})
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.Session{}, nil
	case err != nil:
		return entity.Session{}, err
	}

	return res.ToEntity(), nil
}

func (r *Repository) GetSessionsForAccount(
	ctx context.Context,
	userID uuid.UUID,
	page, size int32,
) (entity.SessionsCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	sessions, err := r.sql.GetSessionsByAccountID(ctx, pgdb.GetSessionsByAccountIDParams{
		AccountID: userID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return entity.SessionsCollection{}, err
	}

	total, err := r.sql.CountSessionsByAccountID(ctx, userID)
	if err != nil {
		return entity.SessionsCollection{}, err
	}

	result := make([]entity.Session, 0, len(sessions))
	for _, s := range sessions {
		result = append(result, s.ToEntity())
	}

	return entity.SessionsCollection{
		Data:  result,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (r *Repository) GetSessionToken(ctx context.Context, sessionID uuid.UUID) (string, error) {
	row, err := r.sql.GetSessionByID(ctx, sessionID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return "", nil
	case err != nil:
		return "", err
	}

	return row.GetHashToken(), nil
}

func (r *Repository) UpdateSessionToken(
	ctx context.Context,
	sessionID uuid.UUID,
	token string,
) (entity.Session, error) {
	ses, err := r.sql.UpdateSessionToken(ctx, pgdb.UpdateSessionTokenParams{
		ID:        sessionID,
		HashToken: token,
	})
	if err != nil {
		return entity.Session{}, err
	}

	return ses.ToEntity(), nil
}

func (r *Repository) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	return r.sql.DeleteSessionByID(ctx, sessionID)
}

func (r *Repository) DeleteSessionsForAccount(ctx context.Context, userID uuid.UUID) error {
	return r.sql.DeleteSessionsByAccountID(ctx, userID)
}

func (r *Repository) DeleteAccountSession(ctx context.Context, userID, sessionID uuid.UUID) error {
	return r.sql.DeleteAccountSession(ctx, pgdb.DeleteAccountSessionParams{
		ID:        sessionID,
		AccountID: userID,
	})
}
