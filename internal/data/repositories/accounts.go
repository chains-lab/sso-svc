package repositories

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/recovery-flow/roles"
	redisrepo "github.com/recovery-flow/sso-oauth/internal/data/dbx/redisdb/repositories"
	sqlrepo "github.com/recovery-flow/sso-oauth/internal/data/dbx/sql/repositories"
	"github.com/recovery-flow/sso-oauth/internal/data/models"
)

type Accounts interface {
	Create(r *http.Request, email string, role roles.UserRole) (*models.Account, error)

	GetByID(r *http.Request, id uuid.UUID) (*models.Account, error)
	GetByEmail(r *http.Request, email string) (*models.Account, error)

	UpdateRole(r *http.Request, id uuid.UUID, role roles.UserRole) (*models.Account, error)
}

type accounts struct {
	redis redisrepo.Accounts
	sql   sqlrepo.Accounts
}

func NewAccounts(redis redisrepo.Accounts, sql sqlrepo.Accounts) Accounts {
	return &accounts{
		redis: redis,
		sql:   sql,
	}
}

func (a *accounts) Create(r *http.Request, email string, role roles.UserRole) (*models.Account, error) {
	acc, err := a.sql.Create(r, email, role)
	if err != nil {
		return nil, err
	}

	res := models.Account{
		ID:        acc.ID,
		Email:     acc.Email,
		Role:      acc.Role,
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	}
	err = a.redis.Add(r.Context(), res, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (a *accounts) GetByEmail(r *http.Request, email string) (*models.Account, error) {
	user, err := a.redis.GetByEmail(r.Context(), email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	acc, err := a.sql.GetByEmail(r, email)
	if err != nil {
		return nil, err
	}
	res := models.Account{
		ID:        acc.ID,
		Email:     acc.Email,
		Role:      acc.Role,
		UpdatedAt: acc.UpdatedAt,
		CreatedAt: acc.CreatedAt,
	}
	err = a.redis.Add(r.Context(), res, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (a *accounts) GetByID(r *http.Request, id uuid.UUID) (*models.Account, error) {
	user, err := a.redis.GetByID(r.Context(), id.String())
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	acc, err := a.sql.GetByID(r, id)
	if err != nil {
		return nil, err
	}
	res := models.Account{
		ID:        acc.ID,
		Email:     acc.Email,
		Role:      acc.Role,
		UpdatedAt: acc.UpdatedAt,
		CreatedAt: acc.CreatedAt,
	}
	err = a.redis.Add(r.Context(), res, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (a *accounts) UpdateRole(r *http.Request, id uuid.UUID, role roles.UserRole) (*models.Account, error) {
	acc, err := a.sql.UpdateRole(r, id, role)
	if err != nil {
		return nil, err
	}

	res := models.Account{
		ID:        acc.ID,
		Email:     acc.Email,
		Role:      acc.Role,
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	}
	err = a.redis.Add(r.Context(), res, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
