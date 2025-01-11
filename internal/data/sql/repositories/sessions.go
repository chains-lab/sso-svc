package repositories

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/sso-oauth/internal/data/sql/repositories/sqlcore"
)

type Sessions interface {
	Create(r *http.Request, userID uuid.UUID, deviceId uuid.UUID, token string) (sqlcore.Session, error)

	GetByID(r *http.Request, id uuid.UUID) (sqlcore.Session, error)
	GetSession(r *http.Request, id uuid.UUID, userID uuid.UUID) (sqlcore.Session, error)
	GetSessions(r *http.Request, userID uuid.UUID) ([]sqlcore.Session, error)

	GetToken(r *http.Request, id uuid.UUID, userID uuid.UUID) (string, error)
	UpdateToken(r *http.Request, id uuid.UUID, token string) error

	DeleteAll(r *http.Request, id uuid.UUID) error
	Delete(r *http.Request, id uuid.UUID, userID uuid.UUID) error

	TerminateSessions(
		r *http.Request,
		userId uuid.UUID,
		curDevId uuid.UUID,
	) error
}

type sessions struct {
	queries *sqlcore.Queries
}

func NewSession(queries *sqlcore.Queries) Sessions {
	return &sessions{queries: queries}
}

func (s *sessions) Create(r *http.Request, userID uuid.UUID, deviceId uuid.UUID, token string) (sqlcore.Session, error) {
	return s.queries.CreateSession(r.Context(), sqlcore.CreateSessionParams{
		ID:      deviceId,
		UserID:  userID,
		Token:   token,
		Client:  httpkit.GetUserAgent(r),
		IpFirst: httpkit.GetClientIP(r),
		IpLast:  httpkit.GetClientIP(r),
	})
}

func (s *sessions) GetByID(r *http.Request, id uuid.UUID) (sqlcore.Session, error) {
	return s.queries.GetSession(r.Context(), id)
}

func (s *sessions) GetSession(r *http.Request, id uuid.UUID, userID uuid.UUID) (sqlcore.Session, error) {
	return s.queries.GetUserSession(r.Context(), sqlcore.GetUserSessionParams{
		ID:     id,
		UserID: userID,
	})
}

func (s *sessions) GetSessions(r *http.Request, userID uuid.UUID) ([]sqlcore.Session, error) {
	return s.queries.GetSessionsByUserID(r.Context(), userID)
}

func (s *sessions) GetToken(r *http.Request, id uuid.UUID, userID uuid.UUID) (string, error) {
	return s.queries.GetSessionToken(r.Context(), sqlcore.GetSessionTokenParams{
		ID:     id,
		UserID: userID,
	})
}

func (s *sessions) UpdateToken(r *http.Request, id uuid.UUID, token string) error {
	return s.queries.UpdateSessionToken(r.Context(), sqlcore.UpdateSessionTokenParams{
		ID:     id,
		Token:  token,
		IpLast: httpkit.GetClientIP(r),
	})
}

func (s *sessions) TerminateSessions(
	r *http.Request,
	userId uuid.UUID,
	curDevId uuid.UUID,
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
		if dev.ID == curDevId {
			continue
		}
		err = queries.DeleteUserSession(ctx, sqlcore.DeleteUserSessionParams{
			ID:     dev.ID,
			UserID: userId,
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *sessions) DeleteAll(r *http.Request, id uuid.UUID) error {
	return s.queries.DeleteUserSessions(r.Context(), id)
}

func (s *sessions) Delete(r *http.Request, id uuid.UUID, userID uuid.UUID) error {
	return s.queries.DeleteUserSession(r.Context(), sqlcore.DeleteUserSessionParams{
		ID:     id,
		UserID: userID,
	})
}
