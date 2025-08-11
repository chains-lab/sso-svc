package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/sso-svc/internal/app/jwtmanager"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/chains-lab/sso-svc/internal/pagination"
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
}

type JWTManager interface {
	EncryptAccess(token string) (string, error)
	EncryptRefresh(token string) (string, error)
	DecryptRefresh(encryptedToken string) (string, error)

	GenerateAccess(
		userID uuid.UUID,
		sessionID uuid.UUID,
		idn string,
	) (string, error)

	GenerateRefresh(
		userID uuid.UUID,
		sessionID uuid.UUID,
		idn string,
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

	refresh, err := s.jwt.GenerateRefresh(user.ID, id, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
	}

	refreshCrypto, err := s.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
	}

	access, err := s.jwt.GenerateAccess(user.ID, id, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
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
			return models.Session{}, models.TokensPair{}, errx.RaiseUserNotFound(
				ctx,
				fmt.Errorf("failed to create session for user %s: %w", user.ID, err),
				user.ID.String(),
			)
		default:
			return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
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

func (s Sessions) Get(ctx context.Context, sessionID, userID uuid.UUID) (models.Session, error) {
	session, err := s.queries.New().FilterID(sessionID).FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Session{}, errx.RaiseSessionNotFound(
				ctx,
				fmt.Errorf("session with id: %s not found for user %s", sessionID, userID),
				sessionID.String(),
				userID.String(),
			)
		default:
			return models.Session{}, errx.RaiseInternal(ctx, err)
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

func (s Sessions) SelectByUserID(ctx context.Context, userID uuid.UUID, pag pagination.Request) ([]models.Session, pagination.Response, error) {
	limit, offset := pagination.CalculateLimitOffset(pag)

	sessions, err := s.queries.New().FilterID(userID).Page(limit, offset).Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, pagination.Response{}, errx.RaiseSessionsForUserNotFound(
				ctx, fmt.Errorf("no sessions found for user %s", userID),
			)
		default:
			return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
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

	return result, pagination.Response{}, nil
}

func (s Sessions) Refresh(ctx context.Context, sessionID uuid.UUID, user models.User, client, token string) (models.Session, models.TokensPair, error) {
	session, err := s.Get(ctx, user.ID, sessionID)
	if err != nil {
		return models.Session{}, models.TokensPair{}, err
	}

	if session.Client != client {
		return models.Session{}, models.TokensPair{}, errx.RaiseSessionClientMismatch(
			ctx,
			fmt.Errorf("client mismatch: expected"),
		)
	}

	access, err := s.jwt.GenerateAccess(session.UserID, session.ID, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
	}

	oldRefresh, err := s.jwt.DecryptRefresh(session.Token)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
	}

	if oldRefresh != token {
		return models.Session{}, models.TokensPair{}, errx.RaiseSessionTokenMismatch(
			ctx,
			fmt.Errorf("refresh token mismatch"),
		)
	}

	newRefresh, err := s.jwt.GenerateRefresh(session.UserID, session.ID, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
	}

	refreshCrypto, err := s.jwt.EncryptRefresh(newRefresh)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
	}

	LastUsed := time.Now().UTC()

	err = s.queries.New().FilterID(sessionID).Update(ctx, dbx.SessionUpdateInput{
		Token:    &refreshCrypto,
		LastUsed: LastUsed,
	})
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
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
			return errx.RaiseUserNotFound(
				ctx,
				fmt.Errorf("no sessions found for user %s", userID),
				userID.String(),
			)
		default:
			return errx.RaiseInternal(ctx, err)
		}
	}
	return nil
}

func (s Sessions) Delete(ctx context.Context, sessionID, userID uuid.UUID) error {
	err := s.queries.New().FilterID(sessionID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.RaiseSessionNotFound(
				ctx,
				fmt.Errorf("session with id: %s not found", sessionID),
				sessionID.String(),
				userID.String(),
			)
		default:
			return errx.RaiseInternal(ctx, err)
		}
	}
	return nil
}
