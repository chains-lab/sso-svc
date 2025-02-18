package entities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/core/models"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/core/tools"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/repository"
	"github.com/recovery-flow/sso-oauth/resources"
	"github.com/sirupsen/logrus"
)

type Session interface {
	Get(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Session, error)

	Delete(ctx context.Context, sessionID uuid.UUID) error
	Terminate(ctx context.Context, userID uuid.UUID, excludeSessionID *uuid.UUID) error

	TokenRefValidate(ctx context.Context, session *models.Session, sk string, token string) error
	TokenUpdate(ctx context.Context, sessionID uuid.UUID, newToken string) error
	TokensGenerate(ctx context.Context, userID uuid.UUID, role string, sessionID uuid.UUID) (*resources.TokensPair, error)
	TokensRefresh(ctx context.Context, accessToken, refreshToken string) (*resources.TokensPair, error)

	AddToBL(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error
	Logout(ctx context.Context, sessionID, userID uuid.UUID) error
}

type session struct {
	Repo repository.Sessions
	Log  *logrus.Logger
}

func NewSession(sessionRepo repository.Sessions, logger *logrus.Logger) Session {
	return &session{
		Repo: sessionRepo,
		Log:  logger,
	}
}

func (s *session) Get(ctx context.Context, sessionID, userID uuid.UUID) (*models.Session, error) {
	ses, err := s.Repo.GetByID(ctx, sessionID)
	if err != nil {
		s.Log.Errorf("Failed to retrieve user session: %v", err)
		return nil, fmt.Errorf("failed to retrieve user session: %w", err)
	}

	if ses.UserID != userID {
		s.Log.Debugf("Session doesn't belong to user")
		return nil, problems.Forbidden("Session doesn't belong to user")
	}

	return ses, nil
}

func (s *session) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Session, error) {
	sessions, err := s.Repo.SelectByUserID(ctx, userID)
	if err != nil {
		s.Log.Errorf("Failed to retrieve user sessions: %v", err)
		return nil, fmt.Errorf("failed to retrieve user sessions: %w", err)
	}

	return sessions, nil
}

func (s *session) Delete(ctx context.Context, sessionID uuid.UUID) error {
	err := s.Repo.Delete(ctx, sessionID)
	if err != nil {
		s.Log.Errorf("Failed to delete session: %v", err)
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (s *session) Terminate(ctx context.Context, userID uuid.UUID, excludeSessionID *uuid.UUID) error {
	err := s.Repo.Terminate(ctx, userID, excludeSessionID)
	if err != nil {
		s.Log.Errorf("Failed to terminate user sessions: %v", err)
		return fmt.Errorf("failed to terminate user sessions: %w", err)
	}

	return nil
}

func (s *session) TokenRefValidate(ctx context.Context, session *models.Session, sk string, token string) error {
	decryptedToken, err := tools.DecryptToken(session.Token, sk)
	if err != nil {
		s.Log.Errorf("Failed to decrypt refresh token: %v", err)
		return err
	}

	if decryptedToken != token {
		s.Log.Warn("Provided refresh token does not match the stored token")
		return err
	}

	return nil
}

func (s *session) TokenUpdate(ctx context.Context, sessionID uuid.UUID, newToken string) error {
	tokenAccess, err := svc.TokenManager.GenerateJWT(
		svc.Config.Server.Name,
		userID.String(),
		svc.Config.JWT.AccessToken.TokenLifetime,
		nil,
		&user.Role,
		&sesIDStr,
	)
	if err != nil {
		svc.Logger.Errorf("Error generating access token: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}
}
