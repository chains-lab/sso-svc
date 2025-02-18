package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/recovery-flow/roles"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/core/models"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/repository/cache"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/repository/sqldb"
	"github.com/redis/go-redis/v9"
)

const (
	ttlAccounts = 15 * time.Minute
)

type Accounts interface {
	Create(ctx context.Context, email string, role roles.UserRole) (*models.Account, error)

	GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error)
	GetByEmail(ctx context.Context, email string) (*models.Account, error)

	UpdateRole(ctx context.Context, id uuid.UUID, role roles.UserRole) (*models.Account, error)
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

func (a *accounts) Create(ctx context.Context, email string, role roles.UserRole) (*models.Account, error) {
	acc, err := a.sql.Create(ctx, email, role)
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
	err = a.redis.Add(ctx, res, ttlAccounts)

	return &res, nil
}

func (a *accounts) GetByEmail(ctx context.Context, email string) (*models.Account, error) {
	user, err := a.redis.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			user = nil
		} else {
			//todo error
		}

	} else if user != nil {
		return user, nil
	}

	acc, err := a.sql.GetByEmail(ctx, email)
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
	err = a.redis.Add(ctx, res, ttlAccounts)
	if err != nil {
		//todo error
	}

	return &res, nil
}

func (a *accounts) GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	user, err := a.redis.GetByID(ctx, id.String())
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

	acc, err := a.sql.GetByID(ctx, id)
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
	err = a.redis.Add(ctx, res, ttlAccounts)
	if err != nil {
		//todo error
	}

	return &res, nil
}

func (a *accounts) UpdateRole(ctx context.Context, id uuid.UUID, role roles.UserRole) (*models.Account, error) {
	acc, err := a.sql.UpdateRole(ctx, id, role)
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
	err = a.redis.Add(ctx, res, ttlAccounts)
	if err != nil {
		//todo error
	}

	return &res, nil
}
