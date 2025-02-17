package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/recovery-flow/sso-oauth/internal/service/data/dbx"
	"github.com/sirupsen/logrus"
)

type Session interface {
	Logout(ctx context.Context, sessionID string, userID string) error
}

type session struct {
	repo   dbx.Sessions
	Logger *logrus.Logger
}

func NewSession(sessionRepo dbx.Sessions, logger *logrus.Logger) Session {
	return &session{
		repo:   sessionRepo,
		Logger: logger,
	}
}

func (s *session) Logout(ctx context.Context, sessionID string, userID string) error {
	log := s.Logger

	sesID, err := uuid.Parse(sessionID)
	if err != nil {
		log.Errorf("Failed to parse session ID: %v", err)
		return fmt.Errorf("failed to parse session ID: %w", err)
	}

	err = s.repo.Delete(ctx, sesID)
	if err != nil {
		log.Errorf("Failed to delete session: %v", err)
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}
