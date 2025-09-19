package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (a App) AdminBlockUser(
	ctx context.Context,
	initiatorID, initiatorSessionID, userID uuid.UUID,
) (models.User, error) {
	_, _, err := a.users.CompareRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return models.User{}, err
	}

	_, err = a.session.GetSessionForInitiator(ctx, initiatorID, initiatorSessionID)
	if err != nil {
		return models.User{}, err
	}

	txErr := a.transaction(func(ctx context.Context) error {
		err := a.session.DeleteAllForUser(ctx, userID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("deleting sessions for user %s: %w", userID, err),
			)
		}

		err = a.users.SetStatus(ctx, userID, enum.UserStatusBlocked)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("deleting user %s: %w", userID, err),
			)
		}

		return nil
	})

	if txErr != nil {
		return models.User{}, txErr
	}

	user, err := a.users.GetByID(ctx, userID)
	if err != nil {
		return models.User{}, txErr
	}

	return user, nil
}

func (a App) AdminUnblockUser(
	ctx context.Context,
	initiatorID, initiatorSessionID, userID uuid.UUID,
) (models.User, error) {
	_, _, err := a.users.CompareRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return models.User{}, err
	}

	_, err = a.session.GetSessionForInitiator(ctx, initiatorID, initiatorSessionID)
	if err != nil {
		return models.User{}, err
	}

	txErr := a.transaction(func(ctx context.Context) error {
		err := a.session.DeleteAllForUser(ctx, userID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("deleting sessions for user %s: %w", userID, err),
			)
		}

		err = a.users.SetStatus(ctx, userID, enum.UserStatusActive)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("deleting user %s: %w", userID, err),
			)
		}

		return nil
	})

	if txErr != nil {
		return models.User{}, txErr
	}

	user, err := a.users.GetByID(ctx, userID)
	if err != nil {
		return models.User{}, txErr
	}

	return user, nil
}
