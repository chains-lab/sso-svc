package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/chains-auth/internal/config"
	"github.com/chains-lab/chains-auth/internal/jwtkit"
	"github.com/chains-lab/chains-auth/internal/repo"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type sessionsRepo interface {
	Create(ctx context.Context, input repo.SessionCreateRequest) error
	Update(ctx context.Context, ID uuid.UUID, input repo.SessionUpdateRequest) error
	Delete(ctx context.Context, ID uuid.UUID) error
	Terminate(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, ID uuid.UUID) (repo.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]repo.Session, error)
	Transaction(fn func(ctx context.Context) error) error
	Drop(ctx context.Context) error
}

type JWTManager interface {
	EncryptAccess(token string) (string, error)
	EncryptRefresh(token string) (string, error)
	DecryptRefresh(encryptedToken string) (string, error)

	GenerateAccess(
		userID uuid.UUID,
		sessionID uuid.UUID,
		subTypeID uuid.UUID,
		idn roles.Role,
	) (string, error)

	GenerateRefresh(
		userID uuid.UUID,
		sessionID uuid.UUID,
		subTypeID uuid.UUID,
		idn roles.Role,
	) (string, error)
}

type Sessions struct {
	repo sessionsRepo
	jwt  JWTManager
}

func NewSession(cfg config.Config, log *logrus.Logger) (Sessions, error) {
	data, err := repo.NewSessions(cfg, log)
	if err != nil {
		return Sessions{}, err
	}

	jwt := jwtkit.NewManager(cfg)

	return Sessions{
		repo: data,
		jwt:  jwt,
	}, nil
}

func (s Sessions) Terminate(ctx context.Context, userUD uuid.UUID) *ape.Error {
	err := s.repo.Terminate(ctx, userUD)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorSessionsForUserNotExist(err)
		default:
			return ape.ErrorInternal(err)
		}
	}
	return nil
}

func (s Sessions) Delete(ctx context.Context, sessionID uuid.UUID) *ape.Error {
	err := s.repo.Delete(ctx, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorSessionDoesNotExist(sessionID, err)
		default:
			return ape.ErrorInternal(err)
		}
	}
	return nil
}

func (s Sessions) Get(ctx context.Context, sessionID uuid.UUID) (models.Session, *ape.Error) {
	session, err := s.repo.GetByID(ctx, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Session{}, ape.ErrorSessionDoesNotExist(sessionID, err)
		default:
			return models.Session{}, ape.ErrorInternal(err)
		}
	}

	return models.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		Token:     session.Token,
		Client:    session.Client,
		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	}, nil
}

func (s Sessions) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Session, *ape.Error) {
	sessions, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ape.ErrorSessionsForUserNotExist(err)
		default:
			return nil, ape.ErrorInternal(err)
		}
	}

	result := make([]models.Session, len(sessions))
	for i, session := range sessions {
		result[i] = models.Session{
			ID:        session.ID,
			UserID:    session.UserID,
			Token:     session.Token,
			Client:    session.Client,
			LastUsed:  session.LastUsed,
			CreatedAt: session.CreatedAt,
		}
	}

	return result, nil
}

func (s Sessions) Create(ctx context.Context, user models.User, client string) (models.Session, models.TokensPair, *ape.Error) {
	id := uuid.New()
	createdAt := time.Now().UTC()

	if user.ID == uuid.Nil {
		return models.Session{}, models.TokensPair{}, ape.ErrorUserDoesNotExist(user.ID, fmt.Errorf("user ID is nil"))
	}

	refresh, err := s.jwt.GenerateRefresh(user.ID, id, user.Subscription, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, ape.ErrorInternal(err)
	}

	refreshCrypto, err := s.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.Session{}, models.TokensPair{}, ape.ErrorInternal(err)
	}

	access, err := s.jwt.GenerateAccess(user.ID, id, user.Subscription, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, ape.ErrorInternal(err)
	}

	err = s.repo.Create(ctx, repo.SessionCreateRequest{
		ID:        id,
		UserID:    user.ID,
		Token:     refreshCrypto,
		Client:    client,
		CreatedAt: createdAt,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Session{}, models.TokensPair{}, ape.ErrorUserDoesNotExist(id, err)
		default:
			return models.Session{}, models.TokensPair{}, ape.ErrorInternal(err)
		}
	}

	return models.Session{
			ID:        id,
			UserID:    user.ID,
			Client:    client,
			Token:     refreshCrypto,
			LastUsed:  createdAt,
			CreatedAt: createdAt,
		}, models.TokensPair{
			Refresh: refresh,
			Access:  access,
		}, nil

}

func (s Sessions) Refresh(ctx context.Context, sessionID uuid.UUID, user models.User, client, token string) (models.Session, models.TokensPair, *ape.Error) {
	session, appErr := s.Get(ctx, sessionID)
	if appErr != nil {
		return models.Session{}, models.TokensPair{}, appErr
	}

	if session.Client != client {
		return models.Session{}, models.TokensPair{}, ape.ErrorSessionClientMismatch(fmt.Errorf("session client mismatch"))
	}

	access, err := s.jwt.GenerateAccess(session.UserID, session.ID, user.Subscription, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, ape.ErrorInternal(err)
	}

	oldRefresh, err := s.jwt.DecryptRefresh(session.Token)
	if err != nil {
		return models.Session{}, models.TokensPair{}, ape.ErrorInternal(err)
	}

	if oldRefresh != token {
		return models.Session{}, models.TokensPair{}, ape.ErrorSessionTokenMismatch(fmt.Errorf("token mismatch"))
	}

	newRefresh, err := s.jwt.GenerateRefresh(session.UserID, session.ID, user.Subscription, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, ape.ErrorInternal(err)
	}

	refreshCrypto, err := s.jwt.EncryptRefresh(newRefresh)
	if err != nil {
		return models.Session{}, models.TokensPair{}, ape.ErrorInternal(err)
	}

	LastUsed := time.Now().UTC()

	err = s.repo.Update(ctx, sessionID, repo.SessionUpdateRequest{
		Token:    &refreshCrypto,
		LastUsed: LastUsed,
	})
	if err != nil {
		return models.Session{}, models.TokensPair{}, ape.ErrorInternal(err)
	}

	return models.Session{
			ID:        session.ID,
			UserID:    session.UserID,
			Client:    session.Client,
			Token:     refreshCrypto,
			LastUsed:  LastUsed,
			CreatedAt: session.CreatedAt,
		}, models.TokensPair{
			Refresh: newRefresh,
			Access:  access,
		}, nil
}
