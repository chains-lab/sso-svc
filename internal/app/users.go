package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (a App) RegisterUser(ctx context.Context, email, password string) error {
	_, err := a.users.GetByEmail(ctx, email)
	if err == nil {
		return errx.ErrorUserAlreadyExists.Raise(
			fmt.Errorf("user with email '%s' already exists", email),
		)
	} else if !errors.Is(err, errx.ErrorUserNotFound) {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("checking existing user with email '%s': %w", email, err),
		)
	}

	err = a.users.Create(ctx, email, password, roles.User)
	if err != nil {
		return err
	}

	return nil
}

func (a App) GetUserByID(ctx context.Context, ID uuid.UUID) (models.User, error) {
	return a.users.GetByID(ctx, ID)
}

func (a App) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	return a.users.GetByEmail(ctx, email)
}

func (a App) UpdatePassword(ctx context.Context, userID, SessionID uuid.UUID, currentPassword, newPassword string) error {
	user, err := a.GetInitiator(ctx, userID, SessionID)
	if err != nil {
		return err
	}

	if err = a.users.CheckPassword(ctx, userID, currentPassword); err != nil {
		return errx.ErrorInvalidLogin.Raise(
			fmt.Errorf("current password is incorrect"),
		)
	}

	err = a.Transaction(func(ctx context.Context) error {
		err = a.users.UpdatePassword(ctx, user.ID, newPassword)
		if err != nil {
			return err
		}

		err = a.session.DeleteAllForUser(ctx, user.ID)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (a App) DeleteOwnUser(ctx context.Context, userID, initiatorSessionID uuid.UUID) error {
	_, err := a.GetInitiator(ctx, userID, initiatorSessionID)
	if err != nil {
		return err
	}

	return a.users.Delete(ctx, userID)
}
