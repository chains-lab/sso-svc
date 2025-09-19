package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (a App) ListOwnSessions(
	ctx context.Context,
	userID uuid.UUID,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Session, pagi.Response, error) {
	return a.session.ListForUser(ctx, userID, pag, sort)
}

func (a App) GetOwnSession(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error) {
	return a.session.GetForUser(ctx, userID, sessionID)
}

func (a App) DeleteOwnSessions(ctx context.Context, userID, initiatorSessionID uuid.UUID) error {
	_, err := a.getInitiator(ctx, userID, initiatorSessionID)
	if err != nil {
		return err
	}

	return a.session.DeleteAllForUser(ctx, userID)
}

func (a App) DeleteOwnSession(ctx context.Context, userID, initiatorSessionID, sessionID uuid.UUID) error {
	_, err := a.getInitiator(ctx, userID, initiatorSessionID)
	if err != nil {
		return err
	}

	return a.session.Delete(ctx, sessionID)
}

func (a App) RefreshSession(
	ctx context.Context,
	token string,
) (models.TokensPair, error) {
	var pair models.TokensPair
	var err error
	txErr := a.transaction(func(ctx context.Context) error {
		pair, err = a.session.Refresh(ctx, token)
		if err != nil {
			return err
		}

		session, err := a.session.Get(ctx, pair.SessionID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get session %s after refresh: %w", pair.SessionID, err),
			)
		}

		_, err = a.users.GetInitiator(ctx, session.UserID)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return models.TokensPair{}, txErr
	}

	return models.TokensPair{
		SessionID: pair.SessionID,
		Refresh:   pair.Refresh,
		Access:    pair.Access,
	}, nil
}
