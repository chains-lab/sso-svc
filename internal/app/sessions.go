package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/constant"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (a App) GoogleLogin(ctx context.Context, email string) (models.TokensPair, error) {
	user, err := a.users.GetByEmail(ctx, email)
	if err != nil {
		return models.TokensPair{}, err
	}

	pair, err := a.session.Create(ctx, user.ID, user.Role, user.EmailVer)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create session for user %s: %w", user.ID, err),
		)
	}

	return pair, nil
}

func (a App) Login(ctx context.Context, email, password string) (models.TokensPair, error) {
	user, err := a.GetUserByEmail(ctx, email)
	if err != nil {
		return models.TokensPair{}, err
	}

	if user.Status == constant.UserStatusBlocked {
		return models.TokensPair{}, errx.ErrorInitiatorIsBlocked.Raise(
			fmt.Errorf("initator user %s is blocked", user.ID),
		)
	}

	if err = a.users.CheckPassword(ctx, user.ID, password); err != nil {
		return models.TokensPair{}, err
	}

	pair, err := a.session.Create(ctx, user.ID, user.Role, user.EmailVer)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create session for user %s: %w", user.ID, err),
		)
	}

	return models.TokensPair{
		SessionID: pair.SessionID,
		Refresh:   pair.Refresh,
		Access:    pair.Access,
	}, nil
}

func (a App) RefreshSessionToken(
	ctx context.Context,
	token string,
) (models.TokensPair, error) {
	var pair models.TokensPair
	var err error
	txErr := a.Transaction(func(ctx context.Context) error {
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

func (a App) GetOwnSession(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error) {
	return a.session.GetForUser(ctx, userID, sessionID)
}

func (a App) SelectOwnSessions(
	ctx context.Context,
	userID uuid.UUID,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Session, pagi.Response, error) {
	return a.session.SelectForUSer(ctx, userID, pag, sort)
}

func (a App) DeleteOwnSession(ctx context.Context, userID, initiatorSessionID, sessionID uuid.UUID) error {
	_, err := a.GetInitiator(ctx, userID, initiatorSessionID)
	if err != nil {
		return err
	}

	return a.session.Delete(ctx, sessionID)
}

func (a App) DeleteOwnSessions(ctx context.Context, userID, initiatorSessionID uuid.UUID) error {
	_, err := a.GetInitiator(ctx, userID, initiatorSessionID)
	if err != nil {
		return err
	}

	return a.session.DeleteAllForUser(ctx, userID)
}
