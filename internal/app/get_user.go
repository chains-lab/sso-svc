package app

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) GetUserByID(ctx context.Context, ID uuid.UUID) (models.User, error) {
	return a.users.GetByID(ctx, ID)
}

func (a App) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	return a.users.GetByEmail(ctx, email)
}
