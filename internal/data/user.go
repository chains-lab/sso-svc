package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/sso-svc/internal/data/pgdb"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (d *Database) CreateUser(ctx context.Context, user models.User, pass models.UserPassword) error {
	err := d.sql.users.New().Insert(ctx, pgdb.User{
		ID:        user.ID,
		Role:      user.SysRole,
		Status:    user.Status,
		Email:     user.Email,
		EmailVer:  user.EmailVer,
		UpdatedAt: user.UpdatedAt,
		CreatedAt: user.CreatedAt,

		PasswordHash: pass.Hash,
		PasswordUpAt: pass.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error) {
	u, err := d.sql.users.New().FilterID(userID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.User{}, nil
	case err != nil:
		return models.User{}, err
	}

	return userSchemaToModel(u), nil
}

func (d *Database) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	u, err := d.sql.users.New().FilterEmail(email).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.User{}, nil
	case err != nil:
		return models.User{}, err
	}

	return userSchemaToModel(u), nil
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

func (d *Database) GetUserPassword(
	ctx context.Context,
	userID uuid.UUID,
) (models.UserPassword, error) {
	u, err := d.sql.users.New().FilterID(userID).Get(ctx)
	if err != nil {
		return models.UserPassword{}, err
	}

	return userSchemaToPasswordDataModel(u), nil
}

func (d *Database) UpdateUserPassword(
	ctx context.Context,
	userID uuid.UUID,
	passwordHash string,
	updatedAt time.Time,
) error {
	return d.sql.users.New().FilterID(userID).UpdatePassword(passwordHash, updatedAt).Update(ctx, updatedAt)
}

func (d *Database) UpdateUserCity(
	ctx context.Context,
	userID uuid.UUID,
	cityID *uuid.UUID,
	cityRole *string,
	updatedAt time.Time,
) error {
	return d.sql.users.New().FilterID(userID).UpdateCity(cityID, cityRole).Update(ctx, updatedAt)
}

func (d *Database) UpdateUserCompany(
	ctx context.Context,
	userID uuid.UUID,
	companyID *uuid.UUID,
	companyRole *string,
	updatedAt time.Time,
) error {
	return d.sql.users.New().FilterID(userID).UpdateCompany(companyID, companyRole).Update(ctx, updatedAt)
}

func (d *Database) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return d.sql.users.New().FilterID(userID).Delete(ctx)
}

func userSchemaToModel(s pgdb.User) models.User {
	res := models.User{
		ID:        s.ID,
		SysRole:   s.Role,
		Status:    s.Status,
		Email:     s.Email,
		EmailVer:  s.EmailVer,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
	if s.CityID.Valid {
		res.CityID = &s.CityID.UUID
	}
	if s.CityRole.Valid {
		res.CityRole = &s.CityRole.String
	}
	if s.CompanyID.Valid {
		res.CompanyID = &s.CompanyID.UUID
	}
	if s.CompanyRole.Valid {
		res.CompanyRole = &s.CompanyRole.String
	}
	return res
}

func userSchemaToPasswordDataModel(s pgdb.User) models.UserPassword {
	return models.UserPassword{
		Hash:      s.PasswordHash,
		UpdatedAt: s.PasswordUpAt,
	}
}
