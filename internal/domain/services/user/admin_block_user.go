package user

import (
	"context"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/sso-svc/internal/data/schemas"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (u Service) AdminBlockUser(ctx context.Context, userID uuid.UUID) (models.User, error) {
	txErr := u.db.Transaction(ctx, func(ctx context.Context) error {
		err := u.db.Sessions().FilterUserID(userID).Delete(ctx)
		if err != nil {
			return err
		}

		err = u.setStatus(ctx, userID, enum.UserStatusBlocked)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return models.User{}, txErr
	}

	user, err := u.GetByID(ctx, userID)
	if err != nil {
		return models.User{}, txErr
	}

	return user, nil
}

func (u Service) AdminUnblockUser(
	ctx context.Context,
	userID uuid.UUID,
) (models.User, error) {
	txErr := u.db.Transaction(ctx, func(ctx context.Context) error {
		err := u.db.Sessions().FilterUserID(userID).Delete(ctx)
		if err != nil {
			return err
		}

		err = u.setStatus(ctx, userID, enum.UserStatusActive)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return models.User{}, txErr
	}

	user, err := u.GetByID(ctx, userID)
	if err != nil {
		return models.User{}, txErr
	}

	return user, nil
}

func (u Service) setStatus(ctx context.Context, userID uuid.UUID, status string) error {
	err := enum.ParseUserStatus(status)
	if err != nil {
		return errx.ErrorUserStatusNotSupported.Raise(
			fmt.Errorf("parsing status for user %s, cause: %w", userID, err),
		)
	}

	err = u.db.Users().FilterID(userID).Update(ctx, schemas.UserUpdateInput{
		Status: &status,
	})
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("updating status for user %s, cause: %w", userID, err),
		)
	}

	return nil
}
