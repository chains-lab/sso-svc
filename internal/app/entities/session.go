package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/app/jwtmanager"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

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

type Session struct {
	query dbx.SessionsQ

	jwt JWTManager
}

func CreateSession(pg *sql.DB, manager jwtmanager.Manager) Session {
	return Session{
		query: dbx.NewSessions(pg),
		jwt:   manager,
	}
}

func (s Session) CreateUserSession(
	ctx context.Context,
	userID uuid.UUID,
	token, client, ip string,
) (models.Session, error) {
	session := dbx.Session{
		ID:        uuid.New(),
		UserID:    userID,
		Client:    client,
		Token:     token,
		IP:        ip,
		LastUsed:  time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}

	err := s.query.Insert(ctx, session)
	if err != nil {
		return models.Session{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create session for user %s: %w", userID, err),
		)
	}

	return models.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		Client:    session.Client,
		IP:        session.IP,
		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	}, nil
}

func (s Session) GetUserSession(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error) {
	session, err := s.query.New().FilterID(sessionID).FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Session{}, errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s not found for user %s", sessionID, userID),
			)
		default:
			return models.Session{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get session with id: %s for user %s: %w", sessionID, userID, err),
			)
		}
	}

	return models.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		Client:    session.Client,
		IP:        session.IP,
		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	}, nil
}

func (s Session) SelectUserSessions(
	ctx context.Context,
	userID uuid.UUID,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Session, pagi.Response, error) {
	if pag.Page == 0 {
		pag.Page = 1
	}
	if pag.Size == 0 {
		pag.Size = 20
	}
	if pag.Size > 100 {
		pag.Size = 100
	}

	limit := pag.Size + 1
	offset := (pag.Page - 1) * pag.Size

	query := s.query.New().Page(limit, offset).FilterUserID(userID)

	for _, sort := range sort {
		ascend := sort.Ascend
		switch sort.Field {
		case "created_at":
			query = query.OrderCreatedAt(ascend)
		default:

		}
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("counting rows: %w", err),
		)
	}

	rows, err := query.Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("selecting rows: %w", err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	result := make([]models.Session, len(rows))
	for i, session := range rows {
		result[i] = models.Session{
			ID:        session.ID,
			UserID:    session.UserID,
			Client:    session.Client,
			IP:        session.IP,
			LastUsed:  session.LastUsed,
			CreatedAt: session.CreatedAt,
		}
	}

	return result, pagi.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: total,
	}, nil
}

func (s Session) DeleteUserSession(ctx context.Context, userID, sessionID uuid.UUID) error {
	err := s.query.New().FilterID(sessionID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s not found", sessionID),
			)
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete session with id: %s for user %s: %w", sessionID, userID, err),
			)
		}
	}
	return nil
}

func (s Session) DeleteUserSessions(ctx context.Context, userID uuid.UUID) error {
	err := s.query.New().FilterUserID(userID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("no sessions found for user %s", userID),
			)
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete sessions for user %s: %w", userID, err),
			)
		}
	}
	return nil
}

func (s Session) UpdateToken(ctx context.Context, userID, sessionID uuid.UUID, role, ip string) (string, error) {
	newRefresh, err := s.jwt.GenerateRefresh(userID, sessionID, role)
	if err != nil {
		return "", errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate refresh token for user %s: %w", userID, err),
		)
	}

	refreshCrypto, err := s.jwt.EncryptRefresh(newRefresh)
	if err != nil {
		return "", errx.ErrorInternal.Raise(
			fmt.Errorf("failed to encrypt refresh token for user %s: %w", userID, err),
		)
	}

	LastUsed := time.Now().UTC()

	err = s.query.New().FilterID(sessionID).Update(ctx, map[string]any{
		"token":     refreshCrypto,
		"ip":        ip,
		"last_used": LastUsed,
	})

	if err != nil {
		return "", err
	}

	return newRefresh, nil
}
