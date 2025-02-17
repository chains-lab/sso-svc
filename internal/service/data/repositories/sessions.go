package repositories

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/data/dbx/cache"
	"github.com/recovery-flow/sso-oauth/internal/service/data/dbx/sqldb"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/models"
	"github.com/redis/go-redis/v9"
)

const (
	ttlSessions = 15 * time.Minute
)

type Sessions interface {
	Create(r *http.Request, userID, sessionID uuid.UUID, token string) (*models.Session, error)

	GetByID(r *http.Request, sessionID uuid.UUID) (*models.Session, error)
	SelectByUserID(r *http.Request, userID uuid.UUID) ([]models.Session, error)

	UpdateToken(r *http.Request, sessionID uuid.UUID, token string) (*models.Session, error)

	Delete(r *http.Request, sessionID uuid.UUID) error

	Terminate(r *http.Request, userID uuid.UUID, sessionID *uuid.UUID) error
}

type sessions struct {
	redis cache.Sessions
	sql   sqldb.Sessions
}

func NewSessions(cfg *config.Config) (Sessions, error) {
	redisRepo := cache.NewSessions(
		redis.NewClient(&redis.Options{
			Addr:     cfg.Database.Redis.Addr,
			Password: cfg.Database.Redis.Password,
			DB:       cfg.Database.Redis.DB,
		}),
		time.Duration(cfg.Database.Redis.Lifetime)*time.Minute,
	)
	sqlRepo, err := sqldb.NewSession(cfg.Database.SQL.URL)
	if err != nil {
		return nil, err
	}
	return &sessions{
		redis: redisRepo,
		sql:   sqlRepo,
	}, nil
}

func (s *sessions) Create(r *http.Request, userID, sessionID uuid.UUID, token string) (*models.Session, error) {
	res, err := s.sql.Create(r, userID, sessionID, token)
	if err != nil {
		return nil, err
	}

	err = s.redis.Add(r.Context(), *res, ttlSessions)
	if err != nil {
		//todo error handling
	}
	return res, nil
}

func (s *sessions) GetByID(r *http.Request, sessionID uuid.UUID) (*models.Session, error) {
	res, err := s.redis.GetByID(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			res = nil
		} else {
			//todo error handling
		}
	} else if res != nil {
		return res, nil
	}

	res, err = s.sql.GetByID(r, sessionID)
	if err != nil {
		return nil, err
	}
	err = s.redis.Add(r.Context(), *res, ttlSessions)
	if err != nil {
		//todo error handling
	}
	return res, nil
}

func (s *sessions) SelectByUserID(r *http.Request, userID uuid.UUID) ([]models.Session, error) {
	res, err := s.redis.SelectByUserID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			res = nil
		} else {
			//todo error handling
		}
	} else if res != nil {
		return res, nil
	}

	res, err = s.sql.SelectByUserID(r, userID)
	if err != nil {
		return nil, err
	}

	go func() {
		for _, ses := range res {
			err = s.redis.Add(r.Context(), ses, ttlSessions)
			if err != nil {
				// todo error handling (лучше логировать ошибку)
			}
		}
	}()

	return res, nil
}

func (s *sessions) UpdateToken(r *http.Request, sessionID uuid.UUID, token string) (*models.Session, error) {
	res, err := s.sql.UpdateToken(r, sessionID, token)
	if err != nil {
		return nil, err
	}

	err = s.redis.Add(r.Context(), *res, ttlSessions)
	if err != nil {
		//todo error handling
	}

	return res, nil
}

func (s *sessions) Delete(r *http.Request, sessionID uuid.UUID) error {
	err := s.sql.Delete(r, sessionID)
	if err != nil {
		return err
	}

	err = s.redis.Delete(r.Context(), sessionID)
	if err != nil {
		//todo error handling
	}

	return nil
}

func (s *sessions) Terminate(r *http.Request, userID uuid.UUID, sessionID *uuid.UUID) error {
	err := s.sql.TerminateSessions(r, userID, sessionID)
	if err != nil {
		return err
	}

	err = s.redis.DeleteByUserID(r.Context(), userID, sessionID)
	if err != nil {
		//todo error handling
	}

	return nil
}
