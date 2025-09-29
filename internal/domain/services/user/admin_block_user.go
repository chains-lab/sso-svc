package user

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/chains-lab/sso-svc/internal/models"
	"github.com/google/uuid"
)

func (s Service) AdminBlockUser(ctx context.Context, userID uuid.UUID) (models.User, error) {
	txErr := s.db.Users().Transaction(ctx, func(ctx context.Context) error {
		err := s.db.Sessions().FilterUserID(userID).Delete(ctx)
		if err != nil {
			return err
		}

		err = s.setStatus(ctx, userID, enum.UserStatusBlocked)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return models.User{}, txErr
	}

	user, err := s.GetByID(ctx, userID)
	if err != nil {
		return models.User{}, txErr
	}

	return user, nil
}

func (s Service) AdminUnblockUser(
	ctx context.Context,
	userID uuid.UUID,
) (models.User, error) {
	txErr := s.db.Users().Transaction(ctx, func(ctx context.Context) error {
		err := s.db.Sessions().FilterUserID(userID).Delete(ctx)
		if err != nil {
			return err
		}

		err = s.setStatus(ctx, userID, enum.UserStatusActive)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return models.User{}, txErr
	}

	user, err := s.GetByID(ctx, userID)
	if err != nil {
		return models.User{}, txErr
	}

	return user, nil
}

func (s Service) setStatus(ctx context.Context, userID uuid.UUID, status string) error {
	err := enum.ParseUserStatus(status)
	if err != nil {
		return errx.ErrorUserStatusNotSupported.Raise(
			fmt.Errorf("parsing status for user %s, cause: %w", userID, err),
		)
	}

	err = s.db.Users().FilterID(userID).UpdateStatus(status).Update(ctx, time.Now().UTC())
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("updating status for user %s, cause: %w", userID, err),
		)
	}

	return nil
}
