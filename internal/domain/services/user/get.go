package user

import (
	"context"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) GetByID(ctx context.Context, ID uuid.UUID) (models.User, error) {
	user, err := s.db.GetUserByID(ctx, ID)
	if err != nil {
		return models.User{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get user with id '%s', cause: %w", ID, err),
		)
	}

	if user == (models.User{}) {
		return models.User{}, errx.ErrorUserNotFound.Raise(
			fmt.Errorf("user with id '%s' not found", ID),
		)
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
	user, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		return models.User{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get user with email '%s', cause: %w", email, err),
		)
	}

	if (user == models.User{}) {
		return models.User{}, errx.ErrorUserNotFound.Raise(
			fmt.Errorf("user with email '%s' not found", email),
		)
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
