package sqldb

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/sso-oauth/internal/service/data/dbx/sqldb/core"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/models"
)

type Sessions interface {
	Create(r *http.Request, userID uuid.UUID, deviceId uuid.UUID, token string) (*models.Session, error)

	GetByID(r *http.Request, id uuid.UUID) (*models.Session, error)
	SelectByUserID(r *http.Request, userID uuid.UUID) ([]models.Session, error)

	UpdateToken(r *http.Request, id uuid.UUID, token string) (*models.Session, error)

	DeleteAll(r *http.Request, userID uuid.UUID) error
	Delete(r *http.Request, id uuid.UUID) error

	TerminateSessions(
		r *http.Request,
		userId uuid.UUID,
		curDevId *uuid.UUID,
	) error
}

type sessions struct {
	queries *core.Queries
}

func NewSession(url string) (Sessions, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &sessions{queries: core.New(db)}, nil
}

func (s *sessions) Create(r *http.Request, userID uuid.UUID, deviceId uuid.UUID, token string) (*models.Session, error) {
	res, err := s.queries.CreateSession(r.Context(), core.CreateSessionParams{
		ID:     deviceId,
		UserID: userID,
		Token:  token,
		Client: httpkit.GetUserAgent(r),
		Ip:     httpkit.GetClientIP(r),
	})
	if err != nil {
		return nil, err
	}

	return parseSession(res), nil
}

func (s *sessions) GetByID(r *http.Request, id uuid.UUID) (*models.Session, error) {
	res, err := s.queries.GetSession(r.Context(), id)
	if err != nil {
		return nil, err
	}

	return parseSession(res), nil
}

func (s *sessions) SelectByUserID(r *http.Request, userID uuid.UUID) ([]models.Session, error) {
	arr, err := s.queries.GetSessionsByUserID(r.Context(), userID)
	if err != nil {
		return nil, err
	}

	var res []models.Session
	for i, session := range arr {
		res[i] = *parseSession(session)
	}

	return res, nil
}

func (s *sessions) UpdateToken(r *http.Request, id uuid.UUID, token string) (*models.Session, error) {
	res, err := s.queries.UpdateSessionToken(r.Context(), core.UpdateSessionTokenParams{
		ID:    id,
		Token: token,
		Ip:    httpkit.GetClientIP(r),
	})
	if err != nil {
		return nil, err
	}

	return parseSession(res), nil
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

func parseSession(session core.Session) *models.Session {
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
