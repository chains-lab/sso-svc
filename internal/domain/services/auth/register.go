package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/data/schemas"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/chains-lab/sso-svc/internal/infra/password"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) Register(
	ctx context.Context,
	email, pass, role string,
) (models.User, error) {
	check, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		return models.User{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get user with email '%s', cause: %w", email, err),
		)
	}

	if (check != schemas.User{}) {
		return models.User{}, errx.ErrorUserAlreadyExists.Raise(
			fmt.Errorf("user with email '%s' not found", email),
		)
	}

	err = roles.ParseRole(role)
	if err != nil {
		return models.User{}, errx.ErrorRoleNotSupported.Raise(
			fmt.Errorf("parsing role for new user with email '%s', cause: %w", email, err),
		)
	}

	err = password.ReliabilityCheck(pass)
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

	err = s.db.CreateUser(ctx, schemas.User{
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

	return models.User{
		ID:        id,
		Role:      role,
		Status:    enum.UserStatusActive,
		Email:     email,
		EmailVer:  false,
		CreatedAt: now,
	}, nil
}

func (s Service) RegisterAdmin(
	ctx context.Context,
	initiatorID uuid.UUID,
	email, pass, role string,
) (models.User, error) {
	initiator, err := s.db.GetUserByID(ctx, initiatorID)
	if err != nil {
		return models.User{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get initiator with id '%s', cause: %w", initiatorID, err),
		)
	}

	if initiator == (schemas.User{}) {
		return models.User{}, errx.ErrorUnauthenticated.Raise(
			fmt.Errorf("initiator with id '%s' not found", initiatorID),
		)
	}

	if initiator.Status == enum.UserStatusBlocked {
		return models.User{}, errx.ErrorInitiatorIsBlocked.Raise(
			fmt.Errorf("user %s is blocked", initiatorID),
		)
	}

	res, err := s.Register(ctx, email, pass, role)
	if err != nil {
		return models.User{}, err
	}

	return res, nil
}
