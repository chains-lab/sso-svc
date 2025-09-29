package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/chains-lab/sso-svc/internal/models"
	"github.com/google/uuid"
)

func (s Service) GetByID(ctx context.Context, ID uuid.UUID) (models.User, error) {
	user, err := s.db.Users().FilterID(ID).Get(ctx)
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

	return models.User{
		ID:        user.ID,
		Role:      user.Role,
		Status:    user.Status,
		Email:     user.Email,
		EmailVer:  user.EmailVer,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s Service) GetByEmail(ctx context.Context, email string) (models.User, error) {
	emailData, err := s.db.Users().FilterEmail(email).Get(ctx)
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

	user, err := s.db.Users().FilterID(emailData.ID).Get(ctx)
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
		Role:      user.Role,
		Status:    user.Status,
		Email:     emailData.Email,
		EmailVer:  emailData.EmailVer,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s Service) GetInitiator(ctx context.Context, userID uuid.UUID) (models.User, error) {
	user, err := s.db.Users().FilterID(userID).Get(ctx)
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

	return models.User{
		ID:        user.ID,
		Role:      user.Role,
		Status:    user.Status,
		Email:     user.Email,
		EmailVer:  user.EmailVer,
		CreatedAt: user.CreatedAt,
	}, nil
}
