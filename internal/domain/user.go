package domain

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/models"
	"github.com/google/uuid"
)

type UserSvc interface {
	AdminBlockUser(ctx context.Context, userID uuid.UUID) (models.User, error)
	AdminUnblockUser(ctx context.Context, userID uuid.UUID) (models.User, error)

	GetByID(ctx context.Context, ID uuid.UUID) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)

	Register(
		ctx context.Context,
		email, pass, role string,
	) (models.User, error)
	RegisterAdmin(
		ctx context.Context,
		initiatorID uuid.UUID,
		email, pass, role string,
	) (models.User, error)

	UpdatePassword(
		ctx context.Context,
		userID uuid.UUID,
		oldPassword, newPassword string,
	) error
}
