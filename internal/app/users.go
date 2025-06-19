package app

import (
	"context"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
)

type usersDomain interface {
	Create(ctx context.Context, email string, role roles.Role) error
	GetByID(ctx context.Context, ID uuid.UUID) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)

	UpdateRole(ctx context.Context, ID uuid.UUID, role roles.Role) error
	UpdateSubscription(ctx context.Context, ID uuid.UUID, subscriptionID uuid.UUID) error
	UpdateVerified(ctx context.Context, ID uuid.UUID, verified bool) error
	UpdateSuspended(ctx context.Context, ID uuid.UUID, suspended bool) error
}

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

func (a App) AdminUpdateUserRole(ctx context.Context, initiatorID, userID uuid.UUID, role roles.Role) error {
	user, err := a.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	initiator, err := a.GetUserByID(ctx, initiatorID)
	if err != nil {
		return err
	}

	if user.Role != roles.SuperUser {
		if roles.CompareRolesUser(initiator.Role, role) < 0 {
			return ape.ErrNoPermission
		}
	}

	if initiator.Suspended {
		return ape.ErrorUserSuspended(initiator.ID)
	}

	err = a.sessions.Terminate(ctx, userID)
	if err != nil {
		return err
	}

	err = a.users.UpdateRole(ctx, userID, role)
	if err != nil {
		return err
	}

	err = a.users.UpdateSubscription(ctx, userID, uuid.Nil)
	if err != nil {
		return err
	}

	err = a.users.UpdateVerified(ctx, userID, true)
	if err != nil {
		return err
	}

	err = a.users.UpdateSuspended(ctx, userID, false)
	if err != nil {
		return err
	}

	err = a.sessions.Terminate(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (a App) AdminUpdateUserSubscription(ctx context.Context, initiatorID, userID, subscriptionID uuid.UUID) error {
	user, err := a.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	initiator, err := a.GetUserByID(ctx, initiatorID)
	if err != nil {
		return err
	}

	if user.Role != roles.SuperUser {
		if roles.CompareRolesUser(initiator.Role, user.Role) < 1 {
			return ape.ErrNoPermission
		}
	}

	if initiator.Suspended {
		return ape.ErrorUserSuspended(initiator.ID)
	}

	err = a.users.UpdateSubscription(ctx, userID, subscriptionID)
	if err != nil {
		return err
	}

	err = a.sessions.Terminate(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (a App) AdminUpdateUserVerified(ctx context.Context, initiatorID, userID uuid.UUID, verified bool) error {
	user, err := a.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	initiator, err := a.GetUserByID(ctx, initiatorID)
	if err != nil {
		return err
	}

	if user.Role != roles.SuperUser {
		if roles.CompareRolesUser(initiator.Role, user.Role) < 1 {
			return ape.ErrNoPermission
		}
	}

	if initiator.Suspended {
		return ape.ErrorUserSuspended(initiator.ID)
	}

	err = a.users.UpdateVerified(ctx, userID, verified)
	if err != nil {
		return err
	}

	err = a.sessions.Terminate(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (a App) AdminUpdateUserSuspended(ctx context.Context, initiatorID, userID uuid.UUID, suspended bool) error {
	user, err := a.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	initiator, err := a.GetUserByID(ctx, initiatorID)
	if err != nil {
		return err
	}

	if user.Role != roles.SuperUser {
		if roles.CompareRolesUser(initiator.Role, user.Role) < 1 {
			return ape.ErrNoPermission
		}
	}

	if initiator.Suspended {
		return ape.ErrorUserSuspended(initiator.ID)
	}

	err = a.users.UpdateSuspended(ctx, userID, suspended)
	if err != nil {
		return err
	}

	err = a.sessions.Terminate(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
