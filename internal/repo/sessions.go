package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/sso-oauth/internal/repo/sqldb"
)

const (
	dataCtxTimeAisle = 10 * time.Second
)

type Session struct {
	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
	Token     string    `json:"token"`
	Client    string    `json:"client"`
	CreatedAt time.Time `json:"created_at"`
	LastUsed  time.Time `json:"last_used"`
}

type SessionSQL interface {
	New() sqldb.SessionsQ
	Insert(ctx context.Context, input sqldb.SessionInsertInput) error
	Update(ctx context.Context, input sqldb.SessionUpdateInput) error
	Delete(ctx context.Context) error
	Count(ctx context.Context) (int, error)
	Select(ctx context.Context) ([]sqldb.SessionModel, error)
	Get(ctx context.Context) (sqldb.SessionModel, error)

	FilterID(id uuid.UUID) sqldb.SessionsQ
	FilterAccountID(accountID uuid.UUID) sqldb.SessionsQ

	Transaction(fn func(ctx context.Context) error) error
	Page(limit, offset uint64) sqldb.SessionsQ
}

type SessionsRepo struct {
	sql sqldb.SessionsQ
}

type SessionCreateRequest struct {
	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
	Token     string    `json:"token"`
	Client    string    `json:"client"`
	LastUsed  time.Time `json:"last_used"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *SessionsRepo) Create(input SessionCreateRequest) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := s.sql.New().Insert(ctxSync, sqldb.SessionInsertInput{
		ID:        input.ID,
		AccountID: input.AccountID,
		Token:     input.Token,
		Client:    input.Client,
		LastUsed:  input.LastUsed,
		CreatedAt: input.CreatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}

type SessionUpdateRequest struct {
	Token    string    `json:"token"`
	Client   string    `json:"client"`
	LastUsed time.Time `json:"last_used"`
}

func (s *SessionsRepo) Update(ID uuid.UUID, input SessionUpdateRequest) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := s.sql.New().FilterID(ID).Update(ctxSync, sqldb.SessionUpdateInput{
		Token:    input.Token,
		Client:   input.Client,
		LastUsed: input.LastUsed,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *SessionsRepo) Delete(ID uuid.UUID) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := s.sql.New().FilterID(ID).Delete(ctxSync)
	if err != nil {
		return err
	}

	return nil
}

func (s *SessionsRepo) Terminate(accountID uuid.UUID) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := s.sql.New().FilterAccountID(accountID).Delete(ctxSync)
	if err != nil {
		return err
	}

	return nil
}

func (s *SessionsRepo) GetByID(ID uuid.UUID) (Session, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	session, err := s.sql.New().FilterID(ID).Get(ctxSync)
	if err != nil {
		return Session{}, err
	}

	return Session{
		ID:        session.ID,
		AccountID: session.AccountID,
		Token:     session.Token,
		Client:    session.Client,
		CreatedAt: session.CreatedAt,
		LastUsed:  session.LastUsed,
	}, nil
}

func (s *SessionsRepo) GetByAccountID(accountID uuid.UUID) ([]Session, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	sessions, err := s.sql.New().FilterAccountID(accountID).Select(ctxSync)
	if err != nil {
		return nil, err
	}

	var result []Session
	for _, session := range sessions {
		result = append(result, Session{
			ID:        session.ID,
			AccountID: session.AccountID,
			Token:     session.Token,
			Client:    session.Client,
			CreatedAt: session.CreatedAt,
			LastUsed:  session.LastUsed,
		})
	}

	return result, nil
}
