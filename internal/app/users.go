package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (a App) RegisterUser(ctx context.Context, email, password string) error {
	_, err := a.GetUserByEmail(ctx, email)
	if err == nil {
		return errx.ErrorUserAlreadyExists.Raise(
			fmt.Errorf("user with email '%s' already exists", email),
		)
	} else if !errors.Is(err, errx.ErrorUserNotFound) {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("checking existing user with email '%s': %w", email, err),
		)
	}

	stmt := dbx.UserModel{
		ID:             uuid.New(),
		Email:          email,
		EmailVer:       false,
		Role:           roles.User,
		EmailUpdatedAt: time.Now().UTC(),
		CreatedAt:      time.Now().UTC(),
	}

	txErr := a.usersQ.Transaction(func(ctx context.Context) error {
		err = a.usersQ.New().Insert(ctx, stmt)
		if err != nil {
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("hashing password for user '%s': %w", stmt.ID, err),
			)
		}

		err = a.passQ.New().Insert(ctx, dbx.UserPasswordModel{
			ID:        stmt.ID,
			PassHash:  string(hash),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

func (a App) GetInitiatorByID(ctx context.Context, ID uuid.UUID) (models.User, error) {
	user, err := a.usersQ.FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, errx.ErrorUnauthenticated.Raise(
				fmt.Errorf("user with id '%s' not found", ID),
			)
		default:
			return models.User{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get user with id '%s': %w", ID, err),
			)
		}
	}

	return userModel(user), nil
}

func (a App) GetUserByID(ctx context.Context, ID uuid.UUID) (models.User, error) {
	user, err := a.usersQ.FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, errx.ErrorUserNotFound.Raise(
				fmt.Errorf("user with id '%s' not found", ID),
			)
		default:
			return models.User{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get user with id '%s': %w", ID, err),
			)
		}
	}

	return userModel(user), nil
}

func (a App) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := a.usersQ.FilterEmail(email).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, errx.ErrorUserNotFound.Raise(
				fmt.Errorf("user with email '%s' not found", email),
			)
		default:
			return models.User{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get user with email '%s': %w", email, err),
			)
		}
	}

	return userModel(user), nil
}

func userModel(model dbx.UserModel) models.User {
	return models.User{
		ID:             model.ID,
		Role:           model.Role,
		Email:          model.Email,
		EmailVer:       model.EmailVer,
		EmailUpdatedAt: model.EmailUpdatedAt,
		CreatedAt:      model.CreatedAt,
	}
}
