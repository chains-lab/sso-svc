package app

import (
	"context"

	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
)

type usersDomain interface {
	Create(ctx context.Context, email string, role roles.Role) error
	UpdateRole(ctx context.Context, ID uuid.UUID, role roles.Role) error
	GetByID(ctx context.Context, ID uuid.UUID) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
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
