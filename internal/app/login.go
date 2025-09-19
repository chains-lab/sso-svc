package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
)

func (a App) Login(ctx context.Context, email, password string) (models.TokensPair, error) {
	user, err := a.GetUserByEmail(ctx, email)
	if err != nil {
		return models.TokensPair{}, err
	}

	if user.Status == enum.UserStatusBlocked {
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

func (a App) LoginByGoogle(ctx context.Context, email string) (models.TokensPair, error) {
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
