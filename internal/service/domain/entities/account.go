package entities

import (
	"context"

	"github.com/google/uuid"
	"github.com/recovery-flow/roles"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/core/models"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/repository"
	"github.com/sirupsen/logrus"
)

type Account interface {
	Get(ctx context.Context, accountID uuid.UUID) (*models.Account, error)
	GetByEmail(ctx context.Context, email string) (*models.Account, error)
	Create(ctx context.Context, acc models.Account) (*models.Account, error)
	UpdateRole(ctx context.Context, accountID uuid.UUID, newRole string) (*models.Account, error)
}

type accounts struct {
	Repo repository.Accounts
	Log  *logrus.Logger
}

func NewAccount(repo repository.Accounts, logger *logrus.Logger) Account {
	return &accounts{
		Repo: repo,
		Log:  logger,
	}
}

func (a *accounts) Get(ctx context.Context, accountID uuid.UUID) (*models.Account, error) {
	user, err := a.Repo.GetByID(ctx, accountID)
	if err != nil {
		a.Log.Errorf("Failed to retrieve user: %v", err)
		return nil, err
	}

	return user, nil
}

func (a *accounts) GetByEmail(ctx context.Context, email string) (*models.Account, error) {
	user, err := a.Repo.GetByEmail(ctx, email)
	if err != nil {
		a.Log.Errorf("Failed to retrieve user: %v", err)
		return nil, err
	}

	return user, nil
}

func (a *accounts) Create(ctx context.Context, account models.Account) (*models.Account, error) {
	accRole, err := roles.StringToRoleUser(account.Role)
	if err != nil {
		return nil, err
	}
	user, err := a.Repo.Create(ctx, account.Email, accRole)
	if err != nil {
		a.Log.Errorf("Failed to create user: %v", err)
		return nil, err
	}

	return user, nil
}

func (a *accounts) UpdateRole(ctx context.Context, accountID uuid.UUID, newRole string) (*models.Account, error) {
	accRole, err := roles.StringToRoleUser(newRole)
	if err != nil {
		return nil, err
	}
	user, err := a.Repo.UpdateRole(ctx, accountID, accRole)
	if err != nil {
		a.Log.Errorf("Failed to update user role: %v", err)
		return nil, err
	}

	return user, nil
}
