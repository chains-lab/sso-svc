package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
)

func (u User) GetByEmail(ctx context.Context, email string) (models.User, error) {
	emailData, err := u.emailQ.New().FilterEmail(email).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, errx.ErrorUserNotFound.Raise(
				fmt.Errorf("user with email '%s' not found, cause: %w", email, err),
			)
		default:
			return models.User{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get user with email '%s', cause: %w", email, err),
			)
		}
	}

	user, err := u.query.FilterID(emailData.ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, errx.ErrorUserNotFound.Raise(
				fmt.Errorf("user with email '%s' not found, cause: %w", email, err),
			)
		default:
			return models.User{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get user with email '%s', cause: %w", email, err),
			)
		}
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
