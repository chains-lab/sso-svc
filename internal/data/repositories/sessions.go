package repositories

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	redisrepo "github.com/recovery-flow/sso-oauth/internal/data/dbx/redisdb/repositories"
	sqlrepo "github.com/recovery-flow/sso-oauth/internal/data/dbx/sql/repositories"
	"github.com/recovery-flow/sso-oauth/internal/data/models"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
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
	redis redisrepo.Sessions
	sql   sqlrepo.Sessions
	log   *logrus.Logger
}

func NewSessions(redis redisrepo.Sessions, sql sqlrepo.Sessions) Sessions {
	return &sessions{
		redis: redis,
		sql:   sql,
	}
}

func (s *sessions) Create(r *http.Request, userID, sessionID uuid.UUID, token string) (*models.Session, error) {
	session, err := s.sql.Create(r, userID, sessionID, token)
	if err != nil {
		return nil, err
	}
	res := models.Session{
		ID:     session.ID,
		UserID: session.UserID,
		Token:  session.Token,
		Client: session.Client,
		IP:     session.Ip,

		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	}

	err = s.redis.Add(r.Context(), res, 15*time.Minute)
	if err != nil {
		s.log.Errorf("error adding session to redis: %v", err)
	}
	return &res, nil
}

func (s *sessions) GetByID(r *http.Request, sessionID uuid.UUID) (*models.Session, error) {
	ses, err := s.redis.GetByID(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			ses = nil
		} else {
			s.log.Errorf("error getting session by id from redis: %v", err)
		}
	} else if ses != nil {
		return ses, nil
	}

	session, err := s.sql.GetByID(r, sessionID)
	if err != nil {
		return nil, err
	}
	res := models.Session{
		ID:     session.ID,
		UserID: session.UserID,
		Token:  session.Token,
		Client: session.Client,
		IP:     session.Ip,

		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	}
	err = s.redis.Add(r.Context(), res, 15*time.Minute)
	if err != nil {
		s.log.Errorf("error adding session to redis: %v", err)
	}
	return &res, nil
}

func (s *sessions) SelectByUserID(r *http.Request, userID uuid.UUID) ([]models.Session, error) {
	ses, err := s.redis.GetByUserID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			ses = nil
		} else {
			s.log.Errorf("error getting session by user id from redis: %v", err)
		}
	} else if ses != nil {
		return ses, nil
	}

	userSessions, err := s.sql.SelectByUserID(r, userID)
	if err != nil {
		return nil, err
	}
	var res []models.Session
	for _, session := range userSessions {
		curSes := models.Session{
			ID:     session.ID,
			UserID: session.UserID,
			Token:  session.Token,
			Client: session.Client,
			IP:     session.Ip,

			LastUsed:  session.LastUsed,
			CreatedAt: session.CreatedAt,
		}
		res = append(res, curSes)

		err = s.redis.Add(r.Context(), curSes, 15*time.Minute)
		if err != nil {
			s.log.Errorf("error adding session to redis: %v", err)
		}
	}

	return res, nil
}

func (s *sessions) UpdateToken(r *http.Request, sessionID uuid.UUID, token string) (*models.Session, error) {
	ses, err := s.sql.UpdateToken(r, sessionID, token)
	if err != nil {
		return nil, err
	}

	res := models.Session{
		ID:        ses.ID,
		UserID:    ses.UserID,
		Token:     ses.Token,
		Client:    ses.Client,
		IP:        ses.Ip,
		CreatedAt: ses.CreatedAt,
		LastUsed:  ses.LastUsed,
	}

	err = s.redis.Add(r.Context(), res, 15*time.Minute)
	if err != nil {
		s.log.Errorf("error adding session to redis: %v", err)
	}

	return &res, nil
}

func (s *sessions) Delete(r *http.Request, sessionID uuid.UUID) error {
	err := s.sql.Delete(r, sessionID)
	if err != nil {
		return err
	}

	err = s.redis.Delete(r.Context(), sessionID)
	if err != nil {
		s.log.Errorf("error deleting session from redis: %v", err)
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
		s.log.Errorf("error deleting session from redis: %v", err)
	}

	return nil
}
