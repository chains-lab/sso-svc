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

func (a App) Register(ctx context.Context, email, password string) error {
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

type RegisterAdminParams struct {
	Email    string
	Password string
	Role     string
}

func (a App) RegisterAdmin(
	ctx context.Context,
	initiatorID, initiatorSessionID uuid.UUID,
	params RegisterAdminParams,
) (models.User, error) {
	_, err := a.users.GetByEmail(ctx, params.Email)
	if !errors.Is(err, errx.ErrorUserNotFound) {
		return models.User{}, err
	}

	initiator, err := a.getInitiator(ctx, initiatorID, initiatorSessionID)
	if err != nil {
		return models.User{}, err
	}

	if initiator.Role == roles.User || initiator.Role == roles.Moder {
		return models.User{}, errx.ErrorNoPermissions.Raise(
			fmt.Errorf("initiator with role %s is not allowed to create user", initiator.Role),
		)
	}

	err = a.users.Create(ctx, params.Email, params.Password, params.Role)
	if err != nil {
		return models.User{}, err
	}

	user, err := a.users.GetByEmail(ctx, params.Email)
	if err != nil {
		return models.User{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get user by email %s after creation: %w", params.Email, err),
		)
	}

	return user, nil
}

func (a App) Register_ONLY_FOR_TESTS(ctx context.Context, email, password, role string) (models.User, error) {
	err := a.users.Create(ctx, email, password, role)
	if err != nil {
		return models.User{}, err
	}

	user, err := a.users.GetByEmail(ctx, email)
	if err != nil {
		return models.User{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get user by email %s after creation: %w", email, err),
		)
	}

	return user, nil
}
