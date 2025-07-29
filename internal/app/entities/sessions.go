package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/ape"
	"github.com/chains-lab/sso-svc/internal/app/jwtmanager"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
)

type sessionsQ interface {
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
	//
	//Drop(ctx context.Context) error
}

type JWTManager interface {
	EncryptAccess(token string) (string, error)
	EncryptRefresh(token string) (string, error)
	DecryptRefresh(encryptedToken string) (string, error)

	GenerateAccess(
		userID uuid.UUID,
		sessionID uuid.UUID,
		idn roles.Role,
	) (string, error)

	GenerateRefresh(
		userID uuid.UUID,
		sessionID uuid.UUID,
		idn roles.Role,
	) (string, error)
}

type Sessions struct {
	queries sessionsQ
	jwt     JWTManager
}

func NewSession(cfg config.Config, log logger.Logger) (Sessions, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return Sessions{}, err
	}

	return Sessions{
		queries: dbx.NewSessions(pg),
		jwt:     jwtmanager.NewManager(cfg),
	}, nil
}

func (s Sessions) Create(ctx context.Context, user models.User, client string) (models.Session, models.TokensPair, error) {
	id := uuid.New()
	createdAt := time.Now().UTC()

	if user.ID == uuid.Nil {
		return models.Session{}, models.TokensPair{}, ape.RaiseUserNotFound(user.ID, fmt.Errorf("user ID is nil"))
	}

	refresh, err := s.jwt.GenerateRefresh(user.ID, id, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, ape.RaiseInternal(err)
	}

	refreshCrypto, err := s.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.Session{}, models.TokensPair{}, ape.RaiseInternal(err)
	}

	access, err := s.jwt.GenerateAccess(user.ID, id, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, ape.RaiseInternal(err)
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
			return models.Session{}, models.TokensPair{}, ape.RaiseUserNotFound(id, err)
		default:
			return models.Session{}, models.TokensPair{}, ape.RaiseInternal(err)
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
			return models.Session{}, ape.RaiseSessionNotFound(sessionID, err)
		default:
			return models.Session{}, ape.RaiseInternal(err)
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
			return nil, ape.RaiseSessionsForUserNotFound(userID, err)
		default:
			return nil, ape.RaiseInternal(err)
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
		return models.Session{}, models.TokensPair{}, ape.RaiseSessionClientMismatch(sessionID, fmt.Errorf("session client mismatch"))
	}

	access, err := s.jwt.GenerateAccess(session.UserID, session.ID, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, ape.RaiseInternal(err)
	}

	oldRefresh, err := s.jwt.DecryptRefresh(session.Token)
	if err != nil {
		return models.Session{}, models.TokensPair{}, ape.RaiseInternal(err)
	}

	if oldRefresh != token {
		return models.Session{}, models.TokensPair{}, ape.RaiseSessionTokenMismatch(sessionID, fmt.Errorf("token mismatch"))
	}

	newRefresh, err := s.jwt.GenerateRefresh(session.UserID, session.ID, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, ape.RaiseInternal(err)
	}

	refreshCrypto, err := s.jwt.EncryptRefresh(newRefresh)
	if err != nil {
		return models.Session{}, models.TokensPair{}, ape.RaiseInternal(err)
	}

	LastUsed := time.Now().UTC()

	err = s.queries.New().FilterID(sessionID).Update(ctx, dbx.SessionUpdateInput{
		Token:    &refreshCrypto,
		LastUsed: LastUsed,
	})
	if err != nil {
		return models.Session{}, models.TokensPair{}, ape.RaiseInternal(err)
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
			return ape.RaiseUserNotFound(userID, err)
		default:
			return ape.RaiseInternal(err)
		}
	}
	return nil
}

func (s Sessions) Delete(ctx context.Context, sessionID uuid.UUID) error {
	err := s.queries.New().FilterID(sessionID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.RaiseSessionNotFound(sessionID, err)
		default:
			return ape.RaiseInternal(err)
		}
	}
	return nil
}
