package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
)

type usersDomain interface {
	Create(ctx context.Context, email string, role roles.Role) *ape.Error
	UpdateRole(ctx context.Context, ID uuid.UUID, role roles.Role) *ape.Error
	GetByID(ctx context.Context, ID uuid.UUID) (models.User, *ape.Error)
	GetByEmail(ctx context.Context, email string) (models.User, *ape.Error)
}

func (a App) UpdateUserRole(ctx context.Context, ID uuid.UUID, role, initiatorRole roles.Role) *ape.Error {
	if roles.CompareRolesUser(role, initiatorRole) != 1 {
		return ape.ErrorUserNoPermissionToUpdateRole(fmt.Errorf("user has no permission to update role"))
	}

	appErr := a.users.UpdateRole(ctx, ID, role)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (a App) GetUserByID(ctx context.Context, ID uuid.UUID) (models.User, *ape.Error) {
	user, appErr := a.users.GetByID(ctx, ID)
	if appErr != nil {
		return models.User{}, appErr
	}

	return user, nil
}

func (a App) GetUserByEmail(ctx context.Context, email string) (models.User, *ape.Error) {
	user, appErr := a.users.GetByEmail(ctx, email)
	if appErr != nil {
		return models.User{}, appErr
	}

	return user, nil
}
