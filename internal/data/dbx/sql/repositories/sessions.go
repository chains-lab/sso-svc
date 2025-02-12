package repositories

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	sqlcore2 "github.com/recovery-flow/sso-oauth/internal/data/dbx/sql/repositories/sqlcore"
)

type Sessions interface {
	Create(r *http.Request, userID uuid.UUID, deviceId uuid.UUID, token string) (sqlcore2.Session, error)

	GetByID(r *http.Request, id uuid.UUID) (sqlcore2.Session, error)
	SelectByUserID(r *http.Request, userID uuid.UUID) ([]sqlcore2.Session, error)

	UpdateToken(r *http.Request, id uuid.UUID, token string) (sqlcore2.Session, error)

	DeleteAll(r *http.Request, userID uuid.UUID) error
	Delete(r *http.Request, id uuid.UUID) error

	TerminateSessions(
		r *http.Request,
		userId uuid.UUID,
		curDevId *uuid.UUID,
	) error
}

type sessions struct {
	queries *sqlcore2.Queries
}

func NewSession(queries *sqlcore2.Queries) Sessions {
	return &sessions{queries: queries}
}

func (s *sessions) Create(r *http.Request, userID uuid.UUID, deviceId uuid.UUID, token string) (sqlcore2.Session, error) {
	return s.queries.CreateSession(r.Context(), sqlcore2.CreateSessionParams{
		ID:     deviceId,
		UserID: userID,
		Token:  token,
		Client: httpkit.GetUserAgent(r),
		Ip:     httpkit.GetClientIP(r),
	})
}

func (s *sessions) GetByID(r *http.Request, id uuid.UUID) (sqlcore2.Session, error) {
	return s.queries.GetSession(r.Context(), id)
}

func (s *sessions) SelectByUserID(r *http.Request, userID uuid.UUID) ([]sqlcore2.Session, error) {
	return s.queries.GetSessionsByUserID(r.Context(), userID)
}

func (s *sessions) UpdateToken(r *http.Request, id uuid.UUID, token string) (sqlcore2.Session, error) {
	return s.queries.UpdateSessionToken(r.Context(), sqlcore2.UpdateSessionTokenParams{
		ID:    id,
		Token: token,
		Ip:    httpkit.GetClientIP(r),
	})
}

func (s *sessions) TerminateSessions(
	r *http.Request,
	userId uuid.UUID,
	curDevId *uuid.UUID,
) error {
	ctx := r.Context()
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

func (s *sessions) DeleteAll(r *http.Request, userID uuid.UUID) error {
	return s.queries.DeleteUserSessions(r.Context(), userID)
}

func (s *sessions) Delete(r *http.Request, id uuid.UUID) error {
	return s.queries.DeleteSession(r.Context(), id)
}
