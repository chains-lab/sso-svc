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

func (a App) Register(ctx context.Context, email, password string) error {
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

type AdminCreateUserInput struct {
	Email    string
	Password string
	Role     string
	Verified bool
}

func (a App) AdminCreateUser(ctx context.Context, initiatorID uuid.UUID, input AdminCreateUserInput) (models.User, error) {
	user, err := a.GetUserByEmail(ctx, input.Email)
	if !errors.Is(err, errx.ErrorUserNotFound) {
		return models.User{}, err
	}

	initiator, err := a.GetUserByID(ctx, initiatorID)
	if err != nil {
		return models.User{}, err
	}

	if initiator.Role == roles.User || initiator.Role == roles.Moder {
		return models.User{}, errx.ErrorNoPermissions.Raise(
			fmt.Errorf("initiator with role %s is not allowed to create user", initiator.Role),
		)
	}

	if user.Role != roles.SuperUser {
		if roles.CompareRolesUser(initiator.Role, input.Role) < 1 {
			return models.User{}, errx.ErrorNoPermissions.Raise(
				fmt.Errorf("initiator Role %s is not allowed to create user Role %s", initiator.Role, input.Role),
			)
		}
	}

	err = a.usersQ.New().Insert(ctx, dbx.UserModel{
		ID:             uuid.New(),
		Email:          input.Email,
		Role:           input.Role,
		EmailVer:       input.Verified,
		EmailUpdatedAt: time.Now().UTC(),
		CreatedAt:      time.Now().UTC(),
	})
	if err != nil {
		return models.User{}, errx.ErrorInternal.Raise(
			fmt.Errorf("creating user with email '%s': %w", input.Email, err),
		)
	}

	user, err = a.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
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
