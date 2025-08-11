package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/app/entity"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (a App) GetUserByID(ctx context.Context, ID uuid.UUID) (models.User, error) {
	user, appErr := a.users.GetByID(ctx, ID)
	if appErr != nil {
		return models.User{}, appErr
	}

	return user, nil
}

func (a App) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	user, appErr := a.users.GetByEmail(ctx, email)
	if appErr != nil {
		return models.User{}, appErr
	}

	return user, nil
}

type AdminCreateUserInput struct {
	Role     string
	Verified bool
}

func (a App) AdminCreateUser(ctx context.Context, initiatorID uuid.UUID, email string, input AdminCreateUserInput) (models.User, error) {
	user, err := a.GetUserByEmail(ctx, email)
	if !errors.Is(err, errx.ErrorUserNotFound) {
		return models.User{}, err
	}

	initiator, err := a.GetUserByID(ctx, initiatorID)
	if err != nil {
		return models.User{}, err
	}

	if initiator.Role == roles.User {
		return models.User{}, errx.RaiseUserRoleIsNotAllowed(
			ctx,
			fmt.Errorf("initiator with role %s is not allowed to create user", initiator.Role),
		)
	}

	if initiator.Suspended {
		return models.User{}, errx.RaiseInitiatorUserSuspended(
			ctx,
			fmt.Errorf("initiator %s is suspended", initiatorID),
			initiatorID.String(),
		)
	}

	if user.Role != roles.SuperUser {
		if roles.CompareRolesUser(initiator.Role, input.Role) < 1 {
			return models.User{}, errx.RaiseInitiatorRoleIsLowThanTarget(
				ctx,
				fmt.Errorf("initiator Role %s is not allowed to create user Role %s", initiator.Role, input.Role),
			)
		}
	}

	err = a.users.Create(ctx, entity.UserCreateInput{
		Email:    email,
		Role:     input.Role,
		Verified: input.Verified,
	})
	if err != nil {
		return models.User{}, errx.RaiseInternal(ctx, err)
	}

	user, err = a.users.GetByEmail(ctx, email)
	if err != nil {
		return models.User{}, errx.RaiseInternal(ctx, err)
	}

	return user, nil
}

func (a App) UpdateUserVerified(ctx context.Context, userID uuid.UUID, verified bool) (models.User, error) {
	// Note: if you want to check initiator rights, not in grpc-service package, with use middleware, you can uncomment the following lines
	//
	//_, user, err := a.users.ComparisonRightsForAdmins(ctx, initiatorID, userID)
	//if err != nil {
	//	return models.User{}, err
	//}

	err := a.users.UpdateVerified(ctx, userID, verified)
	if err != nil {
		return models.User{}, err
	}

	err = a.sessions.Terminate(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	user, err := a.GetUserByID(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (a App) UpdateUserSuspended(ctx context.Context, userID uuid.UUID, suspended bool) (models.User, error) {
	// Note: if you want to check initiator rights, not in grpc-service package, with use middleware, you can uncomment the following lines
	//
	//_, user, err := a.users.ComparisonRightsForAdmins(ctx, initiatorID, userID)
	//if err != nil {
	//	return models.User{}, err
	//}

	err := a.users.UpdateSuspended(ctx, userID, suspended)
	if err != nil {
		return models.User{}, err
	}

	err = a.sessions.Terminate(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	user, err := a.GetUserByID(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
