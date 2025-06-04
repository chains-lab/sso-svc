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
	Terminate(ctx context.Context, accountID uuid.UUID) error
	GetByID(ctx context.Context, ID uuid.UUID) (repo.Session, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]repo.Session, error)
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

func (s Sessions) Terminate(ctx context.Context, accountUD uuid.UUID) *ape.Error {
	err := s.repo.Terminate(ctx, accountUD)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorSessionsForAccountNotExist(err)
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
		AccountID: session.AccountID,
		Client:    session.Client,
		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	}, nil
}

func (s Sessions) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Session, *ape.Error) {
	sessions, err := s.repo.GetByAccountID(ctx, accountID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ape.ErrorSessionsForAccountNotExist(err)
		default:
			return nil, ape.ErrorInternal(err)
		}
	}

	result := make([]models.Session, len(sessions))
	for i, session := range sessions {
		result[i] = models.Session{
			ID:        session.ID,
			AccountID: session.AccountID,
			Client:    session.Client,
			LastUsed:  session.LastUsed,
			CreatedAt: session.CreatedAt,
		}
	}

	return result, nil
}

func (s Sessions) Create(ctx context.Context, account models.Account, client string) (models.Session, *ape.Error) {
	id := uuid.New()
	createdAt := time.Now().UTC()
	token, err := s.jwt.GenerateRefresh(account.ID, id, account.Subscription, account.Role)
	if err != nil {
		return models.Session{}, ape.ErrorInternal(err)
	}

	err = s.repo.Create(ctx, repo.SessionCreateRequest{
		ID:        id,
		AccountID: account.ID,
		Token:     token,
		Client:    client,
		CreatedAt: createdAt,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Session{}, ape.ErrorAccountDoesNotExist(id, err)
		default:
			return models.Session{}, ape.ErrorInternal(err)
		}
	}

	return models.Session{
		ID:        id,
		AccountID: account.ID,
		Client:    client,
		LastUsed:  createdAt,
		CreatedAt: createdAt,
	}, nil

}

func (s Sessions) Refresh(ctx context.Context, sessionID uuid.UUID, account models.Account, client, token string) (models.Session, *ape.Error) {
	session, appErr := s.Get(ctx, sessionID)
	if appErr != nil {
		return models.Session{}, appErr
	}

	if session.Client != client {
		return models.Session{}, ape.ErrorSessionClientMismatch(fmt.Errorf("session client mismatch"))
	}

	access, err := s.jwt.GenerateAccess(session.AccountID, session.ID, account.Subscription, account.Role)
	if err != nil {
		return models.Session{}, ape.ErrorInternal(err)
	}

	refresh, err := s.jwt.GenerateRefresh(session.AccountID, session.ID, account.Subscription, account.Role)
	if err != nil {
		return models.Session{}, ape.ErrorInternal(err)
	}

	refreshCrypto, err := s.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.Session{}, ape.ErrorInternal(err)
	}

	if refreshCrypto != token {
		return models.Session{}, ape.ErrorSessionTokenMismatch(fmt.Errorf("token mismatch"))
	}

	LastUsed := time.Now().UTC()

	err = s.repo.Update(ctx, sessionID, repo.SessionUpdateRequest{
		Token:    &refreshCrypto,
		LastUsed: LastUsed,
	})
	if err != nil {
		return models.Session{}, ape.ErrorInternal(err)
	}

	return models.Session{
		ID:        session.ID,
		AccountID: session.AccountID,
		Access:    access,
		Refresh:   refresh,
		Client:    session.Client,
		LastUsed:  LastUsed,
		CreatedAt: session.CreatedAt,
	}, nil
}
