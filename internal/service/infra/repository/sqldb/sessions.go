package sqldb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/models"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/repository/sqldb/core"
)

type Sessions struct {
	queries *core.Queries
}

func NewSessions(url string) (*Sessions, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Sessions{queries: core.New(db)}, nil
}

func (s *Sessions) Insert(ctx context.Context, session models.Session) (*models.Session, error) {
	res, err := s.queries.CreateSession(ctx, core.CreateSessionParams{
		ID:        session.ID,
		AccountID: session.AccountID,
		Token:     session.Token,
		Client:    session.Client,
		Ip:        session.IP,
	})
	if err != nil {
		return nil, err
	}

	return parseSession(res), nil
}

func (s *Sessions) GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	res, err := s.queries.GetSession(ctx, id)
	if err != nil {
		return nil, err
	}

	return parseSession(res), nil
}

func (s *Sessions) SelectByAccountID(ctx context.Context, AccountID uuid.UUID) ([]models.Session, error) {
	arr, err := s.queries.GetSessionsByAccountID(ctx, AccountID)
	if err != nil {
		return nil, err
	}

	res := make([]models.Session, len(arr))
	for i, session := range arr {
		res[i] = *parseSession(session)
	}

	return res, nil
}

func (s *Sessions) UpdateToken(ctx context.Context, SessionID, AccountID uuid.UUID, token string, IP string) (*models.Session, error) {
	res, err := s.queries.UpdateSessionToken(ctx, core.UpdateSessionTokenParams{
		ID:        SessionID,
		AccountID: AccountID,
		Token:     token,
		Ip:        IP,
	})
	if err != nil {
		return nil, err
	}

	return parseSession(res), nil
}

func (s *Sessions) Terminate(
	ctx context.Context,
	AccountID uuid.UUID,
	curDevId *uuid.UUID,
) error {
	queries, tx, err := s.queries.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = HandleTransactionRollback(tx, err)
	}()

	if err != nil {
		return err
	}

	accountSessions, err := queries.GetSessionsByAccountID(ctx, AccountID)
	if err != nil {
		return err
	}

	for _, dev := range accountSessions {
		if curDevId != nil && dev.ID == *curDevId {
			continue
		}
		err = queries.DeleteSession(ctx, dev.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func HandleTransactionRollback(tx *sql.Tx, originalErr error) error {
	if originalErr != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %v", originalErr, rbErr)
		}
	}
	return originalErr
}

func (s *Sessions) DeleteAll(ctx context.Context, AccountID uuid.UUID) error {
	return s.queries.DeleteAccountSessions(ctx, AccountID)
}

func (s *Sessions) Delete(ctx context.Context, id uuid.UUID) error {
	return s.queries.DeleteSession(ctx, id)
}

func parseSession(session core.Session) *models.Session {
	return &models.Session{
		ID:        session.ID,
		AccountID: session.AccountID,
		Token:     session.Token,
		Client:    session.Client,
		IP:        session.Ip,
		CreatedAt: session.CreatedAt,
		LastUsed:  session.LastUsed,
	}
}
