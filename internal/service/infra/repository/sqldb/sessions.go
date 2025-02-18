package sqldb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/core/models"
	core2 "github.com/recovery-flow/sso-oauth/internal/service/infra/repository/sqldb/core"
)

type Sessions interface {
	Create(ctx context.Context, session models.Session) (*models.Session, error)

	GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error)
	SelectByUserID(ctx context.Context, userID uuid.UUID) ([]models.Session, error)

	UpdateToken(ctx context.Context, id uuid.UUID, token string, IP string) (*models.Session, error)

	DeleteAll(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error

	TerminateSessions(
		ctx context.Context,
		userId uuid.UUID,
		curDevId *uuid.UUID,
	) error
}

type sessions struct {
	queries *core2.Queries
}

func NewSessions(url string) (Sessions, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &sessions{queries: core2.New(db)}, nil
}

func (s *sessions) Create(ctx context.Context, session models.Session) (*models.Session, error) {
	res, err := s.queries.CreateSession(ctx, core2.CreateSessionParams{
		ID:     session.ID,
		UserID: session.UserID,
		Token:  session.Token,
		Client: session.Client,
		Ip:     session.IP,
	})
	if err != nil {
		return nil, err
	}

	return parseSession(res), nil
}

func (s *sessions) GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	res, err := s.queries.GetSession(ctx, id)
	if err != nil {
		return nil, err
	}

	return parseSession(res), nil
}

func (s *sessions) SelectByUserID(ctx context.Context, userID uuid.UUID) ([]models.Session, error) {
	arr, err := s.queries.GetSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var res []models.Session
	for i, session := range arr {
		res[i] = *parseSession(session)
	}

	return res, nil
}

func (s *sessions) UpdateToken(ctx context.Context, id uuid.UUID, token string, IP string) (*models.Session, error) {
	res, err := s.queries.UpdateSessionToken(ctx, core2.UpdateSessionTokenParams{
		ID:    id,
		Token: token,
		Ip:    IP,
	})
	if err != nil {
		return nil, err
	}

	return parseSession(res), nil
}

func (s *sessions) TerminateSessions(
	ctx context.Context,
	userId uuid.UUID,
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

	userSessions, err := queries.GetSessionsByUserID(ctx, userId)
	if err != nil {
		return err
	}

	for _, dev := range userSessions {
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

func (s *sessions) DeleteAll(ctx context.Context, userID uuid.UUID) error {
	return s.queries.DeleteUserSessions(ctx, userID)
}

func (s *sessions) Delete(ctx context.Context, id uuid.UUID) error {
	return s.queries.DeleteSession(ctx, id)
}

func parseSession(session core2.Session) *models.Session {
	return &models.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		Token:     session.Token,
		Client:    session.Client,
		IP:        session.Ip,
		CreatedAt: session.CreatedAt,
		LastUsed:  session.LastUsed,
	}
}
