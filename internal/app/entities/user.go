package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/constant"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	query dbx.UserQ
	passQ dbx.UserPassQ
}

func CreateUser(pg *sql.DB) User {
	return User{
		query: dbx.NewUsers(pg),
	}
}

func (a User) ComparisonRightsForAdmins(
	ctx context.Context,
	initiatorID uuid.UUID,
	userID uuid.UUID,
	dif int,
) (initiator, user models.User, err error) {
	initiator, err = a.GetInitiatorByID(ctx, initiatorID)
	if err != nil {
		return initiator, user, err
	}

	user, err = a.GetUserByID(ctx, userID)
	if err != nil {
		return initiator, user, err
	}

	if user.Role == roles.User {
		return initiator, user, errx.ErrorNoPermissions.Raise(
			fmt.Errorf("user %s is already a user", userID),
		)
	}

	if initiatorID == userID {
		return initiator, user, nil
	}

	if user.Role != roles.SuperUser {
		allowed, err := roles.CompareRolesUser(initiator.Role, user.Role)
		if err != nil {
			return initiator, user, errx.ErrorRoleNotSupported.Raise(
				fmt.Errorf("comparing roles between initiator %s and user %s: %w", initiator.Role, user.Role, err),
			)
		}

		if allowed < dif {
			return initiator, user, errx.ErrorNoPermissions.Raise(
				fmt.Errorf("initiator Role %s is not allowed to interact with this user", initiator.Role),
			)
		}
	}

	return initiator, user, nil
}

func (a User) CheckUserPassword(ctx context.Context, userID uuid.UUID, password string) error {
	secret, err := a.passQ.New().FilterID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.ErrorUserNotFound.Raise(
				fmt.Errorf("password for user %s not found: %w", userID, err),
			)
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("getting password for user %s: %w", userID, err),
			)
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(secret.PassHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("invalid credentials for user %s: %w", userID, err),
			)
		}
		return errx.ErrorInternal.Raise(
			fmt.Errorf("comparing password hash for user %s: %w", userID, err),
		)
	}

	return nil
}

func (a User) CreateUser(ctx context.Context, email, password, role string) error {
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

	err = roles.ParseRole(role)
	if err != nil {
		return errx.ErrorRoleNotSupported.Raise(
			fmt.Errorf("parsing role for new user with email '%s': %w", email, err),
		)
	}

	stmt := dbx.UserModel{
		ID:             uuid.New(),
		Role:           roles.User,
		Status:         constant.UserStatusActive,
		Email:          email,
		EmailVer:       false,
		EmailUpdatedAt: time.Now().UTC(),
		CreatedAt:      time.Now().UTC(),
	}

	txErr := a.query.Transaction(func(ctx context.Context) error {
		err = a.query.New().Insert(ctx, stmt)
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

func (a User) GetInitiatorByID(ctx context.Context, ID uuid.UUID) (models.User, error) {
	user, err := a.query.FilterID(ID).Get(ctx)
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

	if user.Status == constant.UserStatusBlocked {
		return models.User{}, errx.ErrorInitiatorIsBlocked.Raise(
			fmt.Errorf("user with id '%s' is blocked", ID),
		)
	}

	return userModel(user), nil
}

func (a User) GetUserByID(ctx context.Context, ID uuid.UUID) (models.User, error) {
	user, err := a.query.FilterID(ID).Get(ctx)
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

func (a User) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := a.query.FilterEmail(email).Get(ctx)
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

func (a User) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	err := a.passQ.New().FilterID(userID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(err)
	}

	err = a.query.New().FilterID(userID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(err)
	}

	return nil
}

func (a User) SetStatus(ctx context.Context, userID uuid.UUID, status string) error {
	err := constant.ParseUserStatus(status)
	if err != nil {
		return errx.ErrorUserStatusNotSupported.Raise(
			fmt.Errorf("parsing status for user %s: %w", userID, err),
		)
	}

	err = a.query.New().FilterID(userID).Update(ctx,
		map[string]interface{}{"status": status},
	)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("updating status for user %s: %w", userID, err),
		)
	}

	return nil
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
