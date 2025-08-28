package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/gatekit/auth"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/app/jwtmanager"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

type sessionsQ interface {
	New() dbx.SessionsQ
	Insert(ctx context.Context, input dbx.Session) error
	Update(ctx context.Context, input map[string]any) error
	Delete(ctx context.Context) error
	Select(ctx context.Context) ([]dbx.Session, error)
	Get(ctx context.Context) (dbx.Session, error)

	FilterID(id uuid.UUID) dbx.SessionsQ
	FilterUserID(userID uuid.UUID) dbx.SessionsQ

	Page(limit, offset uint64) dbx.SessionsQ
	Count(ctx context.Context) (uint64, error)

	OrderCreatedAt(ascending bool) dbx.SessionsQ
}

type JWTManager interface {
	EncryptAccess(token string) (string, error)
	EncryptRefresh(token string) (string, error)
	DecryptRefresh(encryptedToken string) (string, error)

	ParseRefreshClaims(enc string) (auth.UsersClaims, error)

	GenerateAccess(
		userID uuid.UUID,
		sessionID uuid.UUID,
		idn string,
	) (string, error)

	GenerateRefresh(
		userID uuid.UUID,
		sessionID uuid.UUID,
		role string,
		emailVerified bool,
	) (string, error)
}

type Session struct {
	query sessionsQ

	jwt JWTManager
}

func CreateSession(pg *sql.DB, manager jwtmanager.Manager) Session {
	return Session{
		query: dbx.NewSessions(pg),
		jwt:   manager,
	}
}

func (s Session) Create(
	ctx context.Context,
	userID uuid.UUID,
	role string,
	verified bool,
) (models.TokensPair, error) {
	sessionID := uuid.New()

	access, err := s.jwt.GenerateAccess(userID, sessionID, role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate access token for user %s, cause: %w", userID, err))
	}

	refresh, err := s.jwt.GenerateRefresh(userID, sessionID, role, verified)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate refresh token for user %s, cause: %w", userID, err))
	}

	refreshCrypto, err := s.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to encrypt refresh token for user %s, cause: %w", userID, err))
	}

	session := dbx.Session{
		ID:        sessionID,
		UserID:    userID,
		Token:     refreshCrypto,
		LastUsed:  time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}

	err = s.query.Insert(ctx, session)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create session for user %s, cause: %w", userID, err),
		)
	}

	return models.TokensPair{
		SessionID: sessionID,
		Refresh:   refresh,
		Access:    access,
	}, nil
}

func (s Session) Get(ctx context.Context, sessionID uuid.UUID) (models.Session, error) {
	session, err := s.query.New().FilterID(sessionID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Session{}, errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s not found, cause: %w", sessionID, err),
			)
		default:
			return models.Session{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get session with id: %s cause: %w", sessionID, err),
			)
		}
	}

	return models.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	}, nil
}

func (s Session) GetSessionForInitiator(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error) {
	session, err := s.query.New().FilterUserID(userID).FilterID(sessionID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Session{}, errx.ErrorInitiatorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s not found for user %s, cause: %w", sessionID, userID, err),
			)
		default:
			return models.Session{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get session with id: %s for user %s, cause: %w", sessionID, userID, err),
			)
		}
	}

	return models.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	}, nil
}

func (s Session) GetForUser(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error) {
	session, err := s.query.New().FilterUserID(userID).FilterID(sessionID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Session{}, errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s not found for user %s, cause: %w", sessionID, userID, err),
			)
		default:
			return models.Session{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get session with id: %s for user %s, cause: %w", sessionID, userID, err),
			)
		}
	}

	return models.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	}, nil
}

func (s Session) SelectForUSer(
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
			fmt.Errorf("counting rows, cause: %w", err),
		)
	}

	rows, err := query.Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("selecting rows, cause: %w", err),
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

func (s Session) Delete(ctx context.Context, sessionID uuid.UUID) error {
	err := s.query.New().FilterID(sessionID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s not found, cause: %w", sessionID, err),
			)
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete session with id: %s, cause: %w", sessionID, err),
			)
		}
	}
	return nil
}

func (s Session) DeleteOneForUser(ctx context.Context, userID, sessionID uuid.UUID) error {
	err := s.query.New().FilterUserID(userID).FilterID(sessionID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s not found for user %s, cause: %w", sessionID, userID, err),
			)
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete session with id: %s for user %s, cause: %w", sessionID, userID, err),
			)
		}
	}
	return nil
}

func (s Session) DeleteAllForUser(ctx context.Context, userID uuid.UUID) error {
	err := s.query.New().FilterUserID(userID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.ErrorUserNotFound.Raise(
				fmt.Errorf("no sessions found for user %s, cause: %w", userID, err),
			)
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete sessions for user %s, cause: %w", userID, err),
			)
		}
	}
	return nil
}

func (s Session) Refresh(ctx context.Context, oldRefreshToken string) (models.TokensPair, error) {
	tokenData, err := s.jwt.ParseRefreshClaims(oldRefreshToken)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to decrypt refresh token claims, cause: %w", err),
		)
	}

	userID, err := uuid.Parse(tokenData.Subject)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to parse user id from token claims, cause: %w", err),
		)
	}

	session, err := s.query.New().FilterID(tokenData.Session).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.TokensPair{}, errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s not found for user %s, cause: %w", tokenData.Session, userID, err),
			)
		default:
			return models.TokensPair{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get session with id: %s for user %s, cause: %w", tokenData.Session, userID, err),
			)
		}
	}

	refresh, err := s.jwt.DecryptRefresh(session.Token)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate refresh token for user %s, cause: %w", userID, err),
		)
	}

	if refresh != oldRefreshToken {
		return models.TokensPair{}, errx.ErrorSessionTokenMismatch.Raise(
			fmt.Errorf("refresh token does not match for session %s and user %s, cause: %w", session.ID, userID, err),
		)
	}

	refresh, err = s.jwt.GenerateRefresh(userID, tokenData.Session, tokenData.Role, tokenData.Verified)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate refresh token for user %s, cause: %w", userID, err),
		)
	}

	refreshCrypto, err := s.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to encrypt refresh token for user %s, cause: %w", userID, err),
		)
	}

	access, err := s.jwt.GenerateAccess(userID, tokenData.Session, tokenData.Role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate access token for user %s, cause: %w", userID, err),
		)
	}

	err = s.query.New().FilterID(tokenData.Session).Update(ctx, map[string]any{
		"token":     refreshCrypto,
		"last_used": time.Now().UTC(),
	})

	return models.TokensPair{
		SessionID: tokenData.Session,
		Refresh:   refresh,
		Access:    access,
	}, nil
}
