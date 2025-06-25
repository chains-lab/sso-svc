package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	jwtkit "github.com/chains-lab/chains-auth/internal/app/kit/jwt"
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/chains-auth/internal/dbx"
	"github.com/chains-lab/chains-auth/internal/utils/config"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type SessionsQ interface {
	New() dbx.SessionsQ
	Insert(ctx context.Context, input dbx.SessionModel) error
	Update(ctx context.Context, input dbx.SessionUpdateInput) error
	Delete(ctx context.Context) error
	Count(ctx context.Context) (int, error)
	Select(ctx context.Context) ([]dbx.SessionModel, error)
	Get(ctx context.Context) (dbx.SessionModel, error)

	FilterID(id uuid.UUID) dbx.SessionsQ
	FilterUserID(userID uuid.UUID) dbx.SessionsQ

	Transaction(fn func(ctx context.Context) error) error
	Page(limit, offset uint64) dbx.SessionsQ

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
	queries SessionsQ
	jwt     JWTManager
}

func NewSession(cfg config.Config, log *logrus.Logger) (Sessions, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return Sessions{}, err
	}

	return Sessions{
		queries: dbx.NewSessions(pg),
		jwt:     jwtkit.NewManager(cfg),
	}, nil
}

func (s Sessions) Create(ctx context.Context, user models.User, client string) (models.Session, models.TokensPair, error) {
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

	err = s.queries.New().Insert(ctx, dbx.SessionModel{
		ID:        id,
		UserID:    user.ID,
		Token:     refreshCrypto,
		Client:    client,
		LastUsed:  createdAt,
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

func (s Sessions) Get(ctx context.Context, sessionID uuid.UUID) (models.Session, error) {
	session, err := s.queries.New().FilterID(sessionID).Get(ctx)
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

func (s Sessions) SelectByUserID(ctx context.Context, userID uuid.UUID) ([]models.Session, error) {
	sessions, err := s.queries.New().FilterID(userID).Select(ctx)
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

func (s Sessions) Refresh(ctx context.Context, sessionID uuid.UUID, user models.User, client, token string) (models.Session, models.TokensPair, error) {
	session, err := s.Get(ctx, sessionID)
	if err != nil {
		return models.Session{}, models.TokensPair{}, err
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

	err = s.queries.New().FilterID(sessionID).Update(ctx, dbx.SessionUpdateInput{
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

func (s Sessions) Terminate(ctx context.Context, userID uuid.UUID) error {
	err := s.queries.New().FilterUserID(userID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorUserDoesNotExist(userID, err)
		default:
			return ape.ErrorInternal(err)
		}
	}
	return nil
}

func (s Sessions) Delete(ctx context.Context, sessionID uuid.UUID) error {
	err := s.queries.New().FilterID(sessionID).Delete(ctx)
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
