package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/core/models"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/repository/cache"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/repository/sqldb"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Sessions interface {
	Create(ctx context.Context, session models.Session) (*models.Session, error)

	GetByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	SelectByAccountID(ctx context.Context, userID uuid.UUID) ([]models.Session, error)

	UpdateToken(ctx context.Context, sessionID, UserID uuid.UUID, IP, client, newToken string) (*models.Session, error)

	Delete(ctx context.Context, sessionID uuid.UUID) error

	Terminate(ctx context.Context, userID uuid.UUID, sessionID *uuid.UUID) error
}

type sessions struct {
	redis cache.Sessions
	sql   sqldb.Sessions

	log *logrus.Logger
}

func NewSessions(cfg *config.Config, log *logrus.Logger) (Sessions, error) {
	redisRepo := cache.NewSessions(
		redis.NewClient(&redis.Options{
			Addr:     cfg.Database.Redis.Addr,
			Password: cfg.Database.Redis.Password,
			DB:       cfg.Database.Redis.DB,
		}),
		time.Duration(cfg.Database.Redis.Lifetime)*time.Minute,
	)
	sqlRepo, err := sqldb.NewSessions(cfg.Database.SQL.URL)
	if err != nil {
		return nil, err
	}
	return &sessions{
		redis: redisRepo,
		sql:   *sqlRepo,
		log:   log,
	}, nil
}

func (s *sessions) Create(ctx context.Context, session models.Session) (*models.Session, error) {
	session.CreatedAt = time.Now()
	session.LastUsed = session.CreatedAt
	res, err := s.sql.Insert(ctx, session)
	if err != nil {
		return nil, err
	}

	err = s.redis.Add(ctx, *res)
	if err != nil {
		s.log.WithError(err).Error("error adding session to Redis")
	}
	return res, nil
}

func (s *sessions) GetByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	res, err := s.redis.GetByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			res = nil
		} else {
			s.log.WithError(err).Error("error getting session from Redis")
		}
	} else if res != nil {
		return res, nil
	}

	s.log.WithField("session_id", sessionID)

	res, err = s.sql.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	err = s.redis.Add(ctx, *res)
	if err != nil {
		s.log.WithError(err).Error("error adding session to Redis")
	}

	return res, nil
}

func (s *sessions) SelectByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Session, error) {
	res, err := s.redis.SelectByUserID(ctx, accountID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			res = nil
		} else {
			s.log.WithError(err).Error("error getting session from Redis")
		}
	} else if res != nil {
		return res, nil
	}

	res, err = s.sql.SelectByUserID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	go func() {
		for _, ses := range res {
			err = s.redis.Add(ctx, ses)
			if err != nil {
				s.log.WithError(err).Error("error adding session to Redis")
			}
		}
	}()

	return res, nil
}

func (s *sessions) UpdateToken(ctx context.Context, sessionID, UserID uuid.UUID, IP, client, newToken string) (*models.Session, error) {
	res, err := s.sql.UpdateToken(ctx, sessionID, UserID, newToken, IP)
	if err != nil {
		return nil, err
	}

	err = s.redis.Add(ctx, *res)
	if err != nil {
		s.log.WithError(err).Error("error adding session to Redis")
	}

	return res, nil
}

func (s *sessions) Delete(ctx context.Context, sessionID uuid.UUID) error {
	err := s.sql.Delete(ctx, sessionID)
	if err != nil {
		return err
	}

	err = s.redis.Delete(ctx, sessionID)
	if err != nil {
		s.log.WithError(err).Error("error deleting session from Redis")
	}

	return nil
}

func (s *sessions) Terminate(ctx context.Context, userID uuid.UUID, sessionID *uuid.UUID) error {
	err := s.sql.Terminate(ctx, userID, sessionID)
	if err != nil {
		return err
	}

	err = s.redis.DeleteByUserID(ctx, userID, sessionID)
	if err != nil {
		s.log.WithError(err).Error("error deleting session from Redis")
	}

	return nil
}
