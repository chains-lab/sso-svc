package user

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/sso-svc/internal/domain/enum"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/google/uuid"
)

//TODO this part for the future kafka implementation

// BlockUser - this is methods for lazy block from kafka example nor from http
func (s Service) BlockUser(ctx context.Context, userID uuid.UUID) (models.User, error) {
	_, err := s.GetByID(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	txErr := s.db.Transaction(ctx, func(ctx context.Context) error {
		err = s.db.DeleteAllSessionsForUser(ctx, userID)
		if err != nil {
			return err
		}

		err = s.db.UpdateUserStatus(ctx, userID, enum.UserStatusBlocked, time.Now().UTC())
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("updating status for user %s, cause: %w", userID, err),
			)
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

// UnblockUser - this is methods for lazy unblock from kafka example nor from http
func (s Service) UnblockUser(ctx context.Context, userID uuid.UUID) (models.User, error) {
	_, err := s.GetByID(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	err = s.db.UpdateUserStatus(ctx, userID, enum.UserStatusActive, time.Now().UTC())
	if err != nil {
		return models.User{}, errx.ErrorInternal.Raise(
			fmt.Errorf("updating status for user %s, cause: %w", userID, err),
		)
	}

	user, err := s.GetByID(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
