package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (u User) GetInitiator(ctx context.Context, userID uuid.UUID) (models.User, error) {
	user, err := u.query.FilterID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, errx.ErrorUnauthenticated.Raise(
				fmt.Errorf("user with id '%s' not found, cause: %w", userID, err),
			)
		default:
			return models.User{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get user with id '%s', cause: %w", userID, err),
			)
		}
	}

	emailData, err := u.emailQ.New().FilterID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, errx.ErrorUnauthenticated.Raise(
				fmt.Errorf("email for user with id '%s' not found, cause: %w", userID, err),
			)
		default:
			return models.User{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get email for user with id '%s', cause: %w", userID, err),
			)
		}
	}

	if user.Status == enum.UserStatusBlocked {
		return models.User{}, errx.ErrorInitiatorIsBlocked.Raise(
			fmt.Errorf("user with id '%s' is blocked", userID),
		)
	}

	return models.User{
		ID:        user.ID,
		Email:     emailData.Email,
		Role:      user.Role,
		Status:    user.Status,
		EmailVer:  emailData.Verified,
		CreatedAt: user.CreatedAt,
	}, nil
}
