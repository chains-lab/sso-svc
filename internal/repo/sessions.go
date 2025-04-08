package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/sso-oauth/internal/config"
	"github.com/hs-zavet/sso-oauth/internal/repo/redisdb"
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

type sessionsRedis interface {
	Set(ctx context.Context, input redisdb.SessionCreateInput) error
	Create(ctx context.Context, input redisdb.SessionCreateInput) error
	GetByID(ctx context.Context, sessionID string) (redisdb.SessionModel, error)
	GetByAccountID(ctx context.Context, accountID string) ([]redisdb.SessionModel, error)
	Update(ctx context.Context, sessionID, userID uuid.UUID, update redisdb.SessionUpdateInput) error
	Delete(ctx context.Context, sessionID uuid.UUID) error
	Terminate(ctx context.Context, accountID uuid.UUID) error
	Drop(ctx context.Context) error
}

type SessionsRepo struct {
	sql   sqldb.SessionsQ
	redis sessionsRedis
}

func NewSessions(cfg config.Config) (SessionsRepo, error) {
	db, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return SessionsRepo{}, err
	}

	sessions := sqldb.NewSessions(db)

	return SessionsRepo{
		sql: sessions,
	}, nil
}

type SessionCreateRequest struct {
	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
	Token     string    `json:"token"`
	Client    string    `json:"client"`
	CreatedAt time.Time `json:"created_at"`
}

func (s SessionsRepo) Create(ctx context.Context, input SessionCreateRequest) error {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	err := s.sql.New().Insert(ctxSync, sqldb.SessionInsertInput{
		ID:        input.ID,
		AccountID: input.AccountID,
		Token:     input.Token,
		Client:    input.Client,
		LastUsed:  input.CreatedAt,
		CreatedAt: input.CreatedAt,
	})
	if err != nil {
		return err
	}

	err = s.redis.Create(ctx, redisdb.SessionCreateInput{
		ID:        input.ID,
		AccountID: input.AccountID,
		Token:     input.Token,
		Client:    input.Client,
		LastUsed:  input.CreatedAt,
		CreatedAt: input.CreatedAt,
	})
	if err != nil {
		//Todo log error
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

	session, err := s.sql.New().FilterID(ID).Get(ctxSync)
	if err != nil {
		return err
	}

	err = s.redis.Set(ctx, redisdb.SessionCreateInput{
		ID:        session.ID,
		AccountID: session.AccountID,
		Token:     session.Token,
		Client:    session.Client,
		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	})
	if err != nil {

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

	err = s.redis.Delete(ctx, sessionID)

	return nil
}

func (s SessionsRepo) Terminate(ctx context.Context, accountID uuid.UUID) error {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	err := s.sql.New().FilterAccountID(accountID).Delete(ctxSync)
	if err != nil {
		return err
	}

	err = s.redis.Terminate(ctx, accountID)

	return nil
}

func (s SessionsRepo) GetByID(ctx context.Context, ID uuid.UUID) (Session, error) {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	redisRes, err := s.redis.GetByID(ctx, ID.String())
	if err != nil {
		//TODO: log
	} else {
		return Session{
			ID:        redisRes.ID,
			AccountID: redisRes.AccountID,
			Token:     redisRes.Token,
			Client:    redisRes.Client,
			CreatedAt: redisRes.CreatedAt,
			LastUsed:  redisRes.LastUsed,
		}, nil
	}

	session, err := s.sql.New().FilterID(ID).Get(ctxSync)
	if err != nil {
		return Session{}, err
	}

	if err := s.redis.Set(ctxSync, redisdb.SessionCreateInput{
		ID:        session.ID,
		AccountID: session.AccountID,
		Token:     session.Token,
		Client:    session.Client,
		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	}); err != nil {
		//TODO: log
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

func (s SessionsRepo) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]Session, error) {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	redisRes, err := s.redis.GetByAccountID(ctx, accountID.String())
	if err != nil {
		//TODO: log
	} else {
		var result []Session
		for _, session := range redisRes {
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

		if err := s.redis.Set(ctxSync, redisdb.SessionCreateInput{
			ID:        session.ID,
			AccountID: session.AccountID,
			Token:     session.Token,
			Client:    session.Client,
			CreatedAt: session.CreatedAt,
			LastUsed:  session.LastUsed,
		}); err != nil {
			//TODO: log
		}
	}

	return result, nil
}

func (s SessionsRepo) Transaction(fn func(ctx context.Context) error) error {
	return s.sql.Transaction(fn)
}
