package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/chains-lab/sso-svc/internal/infra/password"
	"github.com/chains-lab/sso-svc/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) Register(
	ctx context.Context,
	email, pass, role string,
) (models.User, error) {
	_, err := s.GetByEmail(ctx, email)
	if err == nil {
		return models.User{}, errx.ErrorUserAlreadyExists.Raise(
			fmt.Errorf("user with email '%s' already exists", email),
		)
	} else if !errors.Is(err, errx.ErrorUserNotFound) {
		return models.User{}, errx.ErrorInternal.Raise(
			fmt.Errorf("checking existing user with email '%s': %w", email, err),
		)
	}

	err = roles.ParseRole(role)
	if err != nil {
		return models.User{}, errx.ErrorRoleNotSupported.Raise(
			fmt.Errorf("parsing role for new user with email '%s', cause: %w", email, err),
		)
	}

	err = password.CheckPassword(pass)
	if err != nil {
		return models.User{}, errx.ErrorPasswordIsInappropriate.Raise(err)
	}

	id := uuid.New()
	now := time.Now().UTC()

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, errx.ErrorInternal.Raise(
			fmt.Errorf("hashing password for user '%s', cause: %w", id, err),
		)
	}

	err = s.db.Users().Insert(ctx, models.UserRow{
		ID:     id,
		Role:   role,
		Status: enum.UserStatusActive,

		PasswordHash: string(hash),
		PasswordUpAt: now,

		Email:    email,
		EmailVer: false,

		UpdatedAt: now,
		CreatedAt: now,
	})
	if err != nil {
		return models.User{}, errx.ErrorInternal.Raise(
			fmt.Errorf("inserting new user with email '%s', cause: %w", email, err),
		)
	}

	user, err := s.GetByID(ctx, id)
	if err != nil {
		return models.User{}, errx.ErrorInternal.Raise(
			fmt.Errorf("getting newly created user with email '%s', cause: %w", email, err),
		)
	}

	return user, nil
}

func (s Service) RegisterAdmin(
	ctx context.Context,
	initiatorID uuid.UUID,
	email, pass, role string,
) (models.User, error) {
	_, err := s.GetByEmail(ctx, email)
	if !errors.Is(err, errx.ErrorUserNotFound) {
		return models.User{}, err
	}

	initiator, err := s.GetInitiator(ctx, initiatorID)
	if err != nil {
		return models.User{}, err
	}

	if initiator.Status == enum.UserStatusBlocked {
		return models.User{}, errx.ErrorInitiatorIsBlocked.Raise(
			fmt.Errorf("initiator %s is blocked", initiator.ID),
		)
	}

	if initiator.Role == roles.User || initiator.Role == roles.Moder {
		return models.User{}, errx.ErrorNoPermissions.Raise(
			fmt.Errorf("initiator with role %s is not allowed to create user", initiator.Role),
		)
	}

	user, err := s.Register(ctx, email, pass, role)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
