package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/core/models"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/repository/cache"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/repository/sqldb"
	"github.com/recovery-flow/tokens/identity"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Accounts interface {
	Create(ctx context.Context, email string, idn identity.IdnType) (*models.Account, error)

	GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error)
	GetByEmail(ctx context.Context, email string) (*models.Account, error)

	UpdateRole(ctx context.Context, id uuid.UUID, idn identity.IdnType) (*models.Account, error)
}

type accounts struct {
	redis cache.Accounts
	sql   sqldb.Accounts
	log   *logrus.Logger
}

func NewAccounts(cfg *config.Config, log *logrus.Logger) (Accounts, error) {
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
		redis: *redisRepo,
		sql:   *sqlRepo,
		log:   log,
	}, nil
}

func (a *accounts) Create(ctx context.Context, email string, idn identity.IdnType) (*models.Account, error) {
	acc, err := a.sql.Insert(ctx, email, idn)
	if err != nil {
		return nil, err
	}
	err = a.redis.Add(ctx, *acc)

	return acc, nil
}

func (a *accounts) GetByEmail(ctx context.Context, email string) (*models.Account, error) {
	user, err := a.redis.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			user = nil
		} else {
			a.log.WithError(err).Error("error getting user from Redis")
		}

	} else if user != nil {
		return user, nil
	}

	acc, err := a.sql.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	err = a.redis.Add(ctx, *acc)
	if err != nil {
		a.log.WithError(err).Error("error adding user to Redis")
	}

	return acc, nil
}

func (a *accounts) GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	user, err := a.redis.GetByID(ctx, id.String())
	if err != nil {
		if errors.Is(err, redis.Nil) {
			user = nil
		} else {
			a.log.WithError(err).Error("error getting user from Redis")
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
	err = a.redis.Add(ctx, *acc)
	if err != nil {
		a.log.WithError(err).Error("error adding user to Redis")
	}

	return acc, nil
}

func (a *accounts) UpdateRole(ctx context.Context, id uuid.UUID, idn identity.IdnType) (*models.Account, error) {
	acc, err := a.sql.UpdateRole(ctx, id, idn)
	if err != nil {
		return nil, err
	}
	err = a.redis.Add(ctx, *acc)
	if err != nil {
		a.log.WithError(err).Error("error adding user to Redis")
	}

	return acc, nil
}
