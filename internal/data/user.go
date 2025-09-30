package data

import (
	"context"
	"time"

	"github.com/chains-lab/sso-svc/internal/data/schemas"
	"github.com/google/uuid"
)

func (d *Database) CreateUser(ctx context.Context, user schemas.User) error {
	return d.sql.users.New().Insert(ctx, user)
}

func (d *Database) GetUserByID(ctx context.Context, userID uuid.UUID) (schemas.User, error) {
	return d.sql.users.New().FilterID(userID).Get(ctx)
}

func (d *Database) GetUserByEmail(ctx context.Context, email string) (schemas.User, error) {
	return d.sql.users.New().FilterEmail(email).Get(ctx)
}

func (d *Database) UpdateUserStatus(
	ctx context.Context,
	userID uuid.UUID,
	status string,
	updatedAt time.Time,
) error {
	return d.sql.users.New().FilterID(userID).UpdateStatus(status).Update(ctx, updatedAt)
}

func (d *Database) UpdateUserEmailVerification(
	ctx context.Context,
	userID uuid.UUID,
	verified bool,
	updatedAt time.Time,
) error {
	return d.sql.users.New().FilterID(userID).UpdateEmailVerified(verified).Update(ctx, updatedAt)
}

func (d *Database) UpdateUserPassword(
	ctx context.Context,
	userID uuid.UUID,
	passwordHash string,
	updatedAt time.Time,
) error {
	return d.sql.users.New().FilterID(userID).UpdatePassword(passwordHash, updatedAt).Update(ctx, updatedAt)
}

func (d *Database) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return d.sql.users.New().FilterID(userID).Delete(ctx)
}
