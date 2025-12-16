package domain

import (
	"context"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) GetOwnSession(ctx context.Context, initiator InitiatorData, sessionID uuid.UUID) (entity.Session, error) {
	_, _, err := s.ValidateSession(ctx, initiator)
	if err != nil {
		return entity.Session{}, err
	}

	session, err := s.db.GetAccountSession(ctx, initiator.AccountID, sessionID)
	if err != nil {
		return entity.Session{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get session with id: %s for account %s, cause: %w", sessionID, initiator.AccountID, err),
		)
	}

	if session.IsNil() {
		return entity.Session{}, errx.ErrorSessionNotFound.Raise(
			fmt.Errorf("session with id: %s for account %s not found", sessionID, initiator.AccountID),
		)
	}

	return session, nil
}

func (s Service) GetOwnSessions(
	ctx context.Context,
	initiator InitiatorData,
	page int32,
	size int32,
) (entity.SessionsCollection, error) {
	_, _, err := s.ValidateSession(ctx, initiator)
	if err != nil {
		return entity.SessionsCollection{}, err
	}

	sessions, err := s.db.GetSessionsForAccount(ctx, initiator.AccountID, page, size)
	if err != nil {
		return entity.SessionsCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to list sessions for account %s, cause: %w", initiator.AccountID, err),
		)
	}

	return sessions, nil
}
