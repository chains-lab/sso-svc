package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (u User) GetByID(ctx context.Context, ID uuid.UUID) (models.User, error) {
	user, err := u.query.FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, errx.ErrorUserNotFound.Raise(
				fmt.Errorf("user with id '%s' not found, cause: %w", ID, err),
			)
		default:
			return models.User{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get user with id '%s', cause: %w", ID, err),
			)
		}
	}

	emailData, err := u.emailQ.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, errx.ErrorUserNotFound.Raise(
				fmt.Errorf("email for user with id '%s' not found, cause: %w", ID, err),
			)
		default:
			return models.User{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get email for user with id '%s', cause: %w", ID, err),
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
