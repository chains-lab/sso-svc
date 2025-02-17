package repositories

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/recovery-flow/roles"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/data/dbx/cache"
	"github.com/recovery-flow/sso-oauth/internal/service/data/dbx/sqldb"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/models"
	"github.com/redis/go-redis/v9"
)

const (
	ttlAccounts = 15 * time.Minute
)

type Accounts interface {
	Create(r *http.Request, email string, role roles.UserRole) (*models.Account, error)

	GetByID(r *http.Request, id uuid.UUID) (*models.Account, error)
	GetByEmail(r *http.Request, email string) (*models.Account, error)

	UpdateRole(r *http.Request, id uuid.UUID, role roles.UserRole) (*models.Account, error)
}

type accounts struct {
	redis cache.Accounts
	sql   sqldb.Accounts
}

func NewAccounts(cfg *config.Config) (Accounts, error) {
	redisRepo := cache.NewAccounts(
		redis.NewClient(&redis.Options{
			Addr:     cfg.Database.Redis.Addr,
			Password: cfg.Database.Redis.Password,
			DB:       cfg.Database.Redis.DB,
		}),
		time.Duration(cfg.Database.Redis.Lifetime)*time.Minute,
	)
	sqlRepo, err := sqldb.NewAccounts(cfg.Database.SQL.URL)
	if err != nil {
		return nil, err
	}
	return &accounts{
		redis: redisRepo,
		sql:   sqlRepo,
	}, nil
}

func (a *accounts) Create(r *http.Request, email string, role roles.UserRole) (*models.Account, error) {
	acc, err := a.sql.Create(r, email, role)
	if err != nil {
		return nil, err
	}

	res := models.Account{
		ID:        acc.ID,
		Email:     acc.Email,
		Role:      acc.Role,
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	}
	err = a.redis.Add(r.Context(), res, ttlAccounts)

	return &res, nil
}

func (a *accounts) GetByEmail(r *http.Request, email string) (*models.Account, error) {
	user, err := a.redis.GetByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			user = nil
		} else {
			//todo error
		}

	} else if user != nil {
		return user, nil
	}

	acc, err := a.sql.GetByEmail(r, email)
	if err != nil {
		return nil, err
	}
	res := models.Account{
		ID:        acc.ID,
		Email:     acc.Email,
		Role:      acc.Role,
		UpdatedAt: acc.UpdatedAt,
		CreatedAt: acc.CreatedAt,
	}
	err = a.redis.Add(r.Context(), res, ttlAccounts)
	if err != nil {
		//todo error
	}

	return &res, nil
}

func (a *accounts) GetByID(r *http.Request, id uuid.UUID) (*models.Account, error) {
	user, err := a.redis.GetByID(r.Context(), id.String())
	if err != nil {
		if errors.Is(err, redis.Nil) {
			user = nil
		} else {
			//todo error
		}
	} else if user != nil {
		return user, nil
	}
	if user != nil {
		return user, nil
	}

	acc, err := a.sql.GetByID(r, id)
	if err != nil {
		return nil, err
	}
	res := models.Account{
		ID:        acc.ID,
		Email:     acc.Email,
		Role:      acc.Role,
		UpdatedAt: acc.UpdatedAt,
		CreatedAt: acc.CreatedAt,
	}
	err = a.redis.Add(r.Context(), res, ttlAccounts)
	if err != nil {
		//todo error
	}

	return &res, nil
}

func (a *accounts) UpdateRole(r *http.Request, id uuid.UUID, role roles.UserRole) (*models.Account, error) {
	acc, err := a.sql.UpdateRole(r, id, role)
	if err != nil {
		return nil, err
	}

	res := models.Account{
		ID:        acc.ID,
		Email:     acc.Email,
		Role:      acc.Role,
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	}
	err = a.redis.Add(r.Context(), res, ttlAccounts)
	if err != nil {
		//todo error
	}

	return &res, nil
}
