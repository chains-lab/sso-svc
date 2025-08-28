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

func (a App) RegisterUser(ctx context.Context, email, password string) error {
	_, err := a.users.GetUserByEmail(ctx, email)
	if err == nil {
		return errx.ErrorUserAlreadyExists.Raise(
			fmt.Errorf("user with email '%s' already exists", email),
		)
	} else if !errors.Is(err, errx.ErrorUserNotFound) {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("checking existing user with email '%s': %w", email, err),
		)
	}

	err = a.users.CreateUser(ctx, email, password, roles.User)
	if err != nil {
		return err
	}

	return nil
}

func (a App) GetUserByID(ctx context.Context, ID uuid.UUID) (models.User, error) {
	return a.users.GetUserByID(ctx, ID)
}

func (a App) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	return a.users.GetUserByEmail(ctx, email)
}
