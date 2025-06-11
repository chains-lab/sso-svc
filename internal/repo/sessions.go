package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/chains-auth/internal/config"
	"github.com/chains-lab/chains-auth/internal/repo/sqldb"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	dataCtxTimeAisle = 10 * time.Second
)

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	Client    string    `json:"client"`
	CreatedAt time.Time `json:"created_at"`
	LastUsed  time.Time `json:"last_used"`
}

type sessionSQL interface {
	New() sqldb.SessionsQ
	Insert(ctx context.Context, input sqldb.SessionInsertInput) error
	Update(ctx context.Context, input sqldb.SessionUpdateInput) error
	Delete(ctx context.Context) error
	Count(ctx context.Context) (int, error)
	Select(ctx context.Context) ([]sqldb.SessionModel, error)
	Get(ctx context.Context) (sqldb.SessionModel, error)

	FilterID(id uuid.UUID) sqldb.SessionsQ
	FilterUserID(userID uuid.UUID) sqldb.SessionsQ

	Transaction(fn func(ctx context.Context) error) error
	Page(limit, offset uint64) sqldb.SessionsQ

	Drop(ctx context.Context) error
}

type SessionsRepo struct {
	sql sessionSQL
}

func NewSessions(cfg config.Config, log *logrus.Logger) (SessionsRepo, error) {
	db, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return SessionsRepo{}, err
	}

	sqlImpl := sqldb.NewSessions(db)

	return SessionsRepo{
		sql: sqlImpl,
	}, nil
}

type SessionCreateRequest struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	Client    string    `json:"client"`
	CreatedAt time.Time `json:"created_at"`
}

func (s SessionsRepo) Create(ctx context.Context, input SessionCreateRequest) error {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	err := s.sql.New().Insert(ctxSync, sqldb.SessionInsertInput{
		ID:        input.ID,
		UserID:    input.UserID,
		Token:     input.Token,
		Client:    input.Client,
		LastUsed:  input.CreatedAt,
		CreatedAt: input.CreatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}

type SessionUpdateRequest struct {
	Token    *string   `json:"token"`
	LastUsed time.Time `json:"last_used"`
}

func (s SessionsRepo) Update(ctx context.Context, ID uuid.UUID, input SessionUpdateRequest) error {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	var sqlInput sqldb.SessionUpdateInput
	if input.Token != nil {
		sqlInput.Token = input.Token
	}
	sqlInput.LastUsed = input.LastUsed

	err := s.sql.New().FilterID(ID).Update(ctxSync, sqldb.SessionUpdateInput{
		Token:    input.Token,
		LastUsed: input.LastUsed,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s SessionsRepo) Delete(ctx context.Context, sessionID uuid.UUID) error {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	err := s.sql.New().FilterID(sessionID).Delete(ctxSync)
	if err != nil {
		return err
	}

	return nil
}

func (s SessionsRepo) Terminate(ctx context.Context, userID uuid.UUID) error {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	err := s.sql.New().FilterUserID(userID).Delete(ctxSync)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil
		default:
			return err
		}
	}

	return nil
}

func (s SessionsRepo) GetByID(ctx context.Context, ID uuid.UUID) (Session, error) {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	session, err := s.sql.New().FilterID(ID).Get(ctxSync)
	if err != nil {
		return Session{}, err
	}

	return Session{
		ID:        session.ID,
		UserID:    session.UserID,
		Token:     session.Token,
		Client:    session.Client,
		CreatedAt: session.CreatedAt,
		LastUsed:  session.LastUsed,
	}, nil
}

func (s SessionsRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]Session, error) {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	sessions, err := s.sql.New().FilterUserID(userID).Select(ctxSync)
	if err != nil {
		return nil, err
	}

	var result []Session
	for _, session := range sessions {
		result = append(result, Session{
			ID:        session.ID,
			UserID:    session.UserID,
			Token:     session.Token,
			Client:    session.Client,
			CreatedAt: session.CreatedAt,
			LastUsed:  session.LastUsed,
		})
	}

	return result, nil
}

func (s SessionsRepo) Transaction(fn func(ctx context.Context) error) error {
	return s.sql.Transaction(fn)
}

func (s SessionsRepo) Drop(ctx context.Context) error {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	if err := s.sql.Drop(ctxSync); err != nil {
		return err
	}

	return nil
}
